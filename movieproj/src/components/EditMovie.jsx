import { useEffect, useState } from "react"
import { useNavigate, useOutletContext, useParams } from "react-router-dom"
import "./Editmovie.css"

import Swal from "sweetalert2"

export default function EditMovie() {
    const { jwttok } = useOutletContext()
const[error,seterror]=useState([])
    const [movie, setMovie] = useState({
        id: 0,
        title: "",
        release_date: "",
        runtime: "",
        mpaa_rating: "",
        description: "",
        genres:[],
        genres_array:[Array(13).fill(false)],
    })



    let { id } = useParams()
    if (id === undefined) {
        id = 0
    }
    const navigate = useNavigate()

    useEffect(() => {
        if (jwttok === "") {
            navigate("/login")
            return
        }

         if(id===0){
          setMovie({
            id: 0,
            title: "",
            release_date: "",
            runtime: "",
            mpaa_rating: "",
            description: "",
            genres:[],
            genres_array:[Array(13).fill(false)],
          })

           const header=new Headers()
           header.append("Content-Type","application/json")
           const requestoptions={
            method:"GET",
            headers:header,
           }
           fetch(`http://localhost:8080/genres`,requestoptions)
           .then(response=>response.json())
           .then((data)=>{
            const checks=[]
            data.forEach((genre)=>{
                checks.push({id:genre.id,checked:false,genre:genre.genre})
            })

             setMovie(m=>({
                ...m,
                genres:checks,
                genres_array:[]
             }))


           })
           .catch(error=>console.error(error))
         }else{
 
            const headers=new Headers()
            headers.append("Content-Type","application/json")
            headers.append("Authorization","Bearer "+jwttok)
            const requestoptions={
                method:"GET",
                headers:headers,

            }
            fetch(`http://localhost:8080/admin/movies/${id}`,requestoptions)
            .then(response=>{
                if(response.status!==200){
                    console.log(err)
                }
                return response.json()
            })
            .then((data)=>{
                data.movie.release_date=new Date(data.movie.release_date).toISOString().split('T')[0]
                const checks=[]
                data.genres.forEach(g=>{
                    if(data.movie.genres_array.indexOf(g.id)!==-1){
                        checks.push({id:g.id,checked:true,genre:g.genre})
                    }
                    else{
                        checks.push({id:g.id,checked:false,genre:g.genre})
                    }
                })
                setMovie({
                    ...data.movie,
                    genres:checks,
                })
            })
            .catch(error=>console.error(error))
         }



    }, [id,jwttok, navigate])


const confirmdelmovie=()=>{
    Swal.fire({
        title: "Delete Movie",
        text: "Are you sure you want to delete the movie",
        icon: "warning",
        showCancelButton: true,
        confirmButtonColor: "#3085d6",
        cancelButtonColor: "#d33",
        confirmButtonText: "Yes, delete it!"
      }).then((result) => {
        if (result.isConfirmed) {
            let headers=new Headers()
            headers.append("Authorization","Bearer "+jwttok)
            const requestoptions={
                method:"DELETE",
                headers:headers,
            }
         fetch(`http://localhost:8080/admin/movies/${movie.id}`,requestoptions)
         .then((response)=>
            response.json()
         )
         .then((data)=>{
            if (data.error){
                console.log(data.error)
            }else{
                navigate("/manage")
            }
         })
         .catch(error=>console.error(error))
        }
      });
      
}



    const handleSubmit = (e) => {
        e.preventDefault();

        let errors=[];
        let required=[
            {field:movie.title,name:"title"},
            {field:movie.release_date,name:"release_date"},
            {field:movie.runtime,name:"runtime"},
            {field:movie.description,name:"description"},
            {field:movie.mpaa_rating,name:"mpaa_rating"},
        ]

        required.forEach(function (obj){
            if(obj.field==""||obj.field==null){
                errors.push(obj.name)
            }
        })

        if(movie.genres_array.length===0){
           Swal.fire({
            title: 'Error',
            text: 'Please select at least one genre',
            icon: 'error',
            confirmButtonText: 'OK',
           })
            errors.push("genres")
        }



     seterror(errors)
     if(errors.length>0){
        return false
     }




    const headers=new Headers()
    headers.append('Authorization', `Bearer ${jwttok}`)
    headers.append('Content-Type', 'application/json')

     let method="PUT"
     if(movie.id>0){
        method="PATCH"
     }
     const requestBody=movie;
     requestBody.release_date=new Date(movie.release_date)
     requestBody.runtime=parseInt(requestBody.runtime)

     let requestoptions={
        body:JSON.stringify(requestBody),
        method:method,
        headers:headers,
        credentials:"include"
     }
fetch(`http://localhost:8080/admin/movies/${movie.id}`,requestoptions)
.then(response => response.json())
.then(data => {
    if(data.error){
        console.log(data.error)
    }
    else{
        navigate("/manage")
    }
})
.catch(err => console.error("Fetch error:", err));

        
    }

    const handleChange = (e) => {
        let name = e.target.name
        let value = e.target.value
        setMovie({
            ...movie,
            [name]: value,
        })
    }
    const handleCheck=(e,position)=>{
      

        let tmparr=movie.genres
        tmparr[position].checked=!tmparr[position].checked
        let tempids=movie.genres_array
        if (!e.target.checked){
            tempids.splice(tempids.indexOf(e.target.value))
        }
        else{
            tempids.push(parseInt(e.target.value,10))
        }
setMovie({
    ...movie,
    genres_array:tempids
})

    }

    return (
        <div className="edit-movie-container">
            <h1>ADD/EDIT Movie</h1>
            {error.length > 0 && (
    <div className="alert alert-danger">
        Invalid form:
        <ul>
        {error.map((err, index) => (
            <li key={index}>
                {err} is required
            </li>
        ))}
        </ul>
    </div>
)}
            <hr />
            <form className="movie-form" onSubmit={handleSubmit}>
                <input type="hidden" name="id" value={movie.id} />
                <div className="form-group">
                    <label htmlFor="title">Title:</label>
                    <input
                        type="text"
                        name="title"
                        value={movie.title}
                        onChange={handleChange}
                        placeholder="Enter movie title"
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="release_date">Release Date:</label>
                    <input
                        type="date"
                        name="release_date"
                        value={movie.release_date}
                        onChange={handleChange}
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="runtime">Runtime:</label>
                    <input
                        type="text"
                        name="runtime"
                        value={movie.runtime}
                        onChange={handleChange}
                        placeholder="Enter runtime"
                    />
                </div>
                <div className="form-group">
                    <label htmlFor="mpaa_rating">MPAA Rating:</label>
                    <select
                        name="mpaa_rating"
                        id="mpaarating"
                        value={movie.mpaa_rating}
                        onChange={handleChange}
                    >
                        <option value="">Choose...</option>
                        <option value="G">G</option>
                        <option value="PG-17">PG-17</option>
                        <option value="PG-13">PG-13</option>
                        <option value="R">R</option>
                        <option value="NC17">NC17</option>
                        <option value="18A">18A</option>
                    </select>
                </div>
                <div className="form-group">
                    <label htmlFor="description">Description:</label>
                    <textarea
                        name="description"
                        value={movie.description}
                        onChange={handleChange}
                        placeholder="Enter movie description"
                    />
                </div>
                <hr />
                <h3>Genre</h3>
                {movie.genres.length > 0 && (
                    <div className="form-group">
                        {movie.genres.map((g,index) => (
                            <div key={g.id}>
                                <input
                                    title={g.genre}
                                    name="genre"
                                    key={index}
                                    id={"genre-"+index}
                                    type="checkbox"
                                    checked={movie.genres[index].checked}
                                    onChange={(e) => handleCheck(e,index)}
                                    value={g.id}
                                />
                                <label>{g.genre}</label>
                            </div>
                        ))}
                    </div>
                )}
                <button className="btn btn-primary">Save</button>

                {movie.id>0&&
                <a href="#!" className="btn btn-danger ms-2" onClick={confirmdelmovie}>Delete Movie</a>}
            </form>
        </div>
    )
}
