
import { useEffect } from "react"
import { useState } from "react"
import { Link } from "react-router-dom"

export default function Graphql() {


    const [movies, setmovies] = useState([])
    const [search, setSearch] = useState('')
    const [list, setlist] = useState([])

    const performsearch = () => {
        const payload = `
   {
   search(titlecontains:"${search}") {
       id
       title
       runtime
       release_date
       mpaa_rating
   }
   }
   
   `;

   const headers=new Headers()
   headers.append('Content-Type', 'application/graphql')



   const requestoption={
    method: 'POST',
    body: payload,
    headers: headers,
    
   }
   fetch(`http://localhost:8080/graph`,requestoption)
   .then(response => response.json())
   .then((response)=>{
    let thelist=Object.values(response.data.search)
    setmovies(thelist)
   })
   .catch((error) => {
    console.error('Error:', error);
   })
    }

    const handlechange = (event) => {
        event.preventDefault()
        let value = event.target.value
        setSearch(value)
        if (value.length > 1) {
            performsearch()
        }
        else {
            setmovies(list)
        }
    }


    useEffect(() => {

        const payload = `
        {
          list{
           id
           title
           release_date
           mpaa_rating
          }
        }
        `;
        const headers = new Headers()
        headers.append('Content-Type', 'application/graphql')
        const requestoption = {
            method: 'POST',
            headers: headers,
            body: payload
        }

        fetch(`http://localhost:8080/graph`, requestoption)
            .then(response => response.json())
            .then((response) => {
                let thelist = Object.values(response.data.list)
                setmovies(thelist)
                setlist(thelist)
            })
            .catch(error => console.error('Error:', error))
    }, [])



    return (
        <>
            <h1>Graphql</h1>
            <hr />

            <form action="" onSubmit={handlechange}>
                <input
                    title="Search"
                    type="search"
                    name="search"
                    value={search}
                    className="form-control"
                    onChange={handlechange}>

                </input>
            </form>
            <hr />
            {movies ? (
                <table className="table table-striped table-hover">
                    <thead>
                        <tr>
                            <th>Movie</th>
                            <th>Release Date</th>
                            <th>Rating</th>
                        </tr>
                    </thead>
                    <tbody>
                        {movies.map((m) => (
                            <tr key={m.id}>
                                <td>
                                    <Link to={`/movies/${m.id}`}>{m.title}</Link>
                                </td>
                                <td>{new Date(m.release_date).toLocaleDateString()}</td>
                                <td>{m.mpaa_rating}</td>
                            </tr>
                        ))}
                    </tbody>
                </table>
            ) : (
                <p>No movies yet</p>
            )}
        </>
    )
}