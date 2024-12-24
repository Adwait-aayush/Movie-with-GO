import { useEffect, useState } from "react"
import { Link } from "react-router-dom";
export default function Movies() {
  const [movies, setmovies] = useState([])

  useEffect(() => {
    const header = new Headers();
    header.append('Content-Type', 'application/json');
    const requestoptions = {
      method: 'GET',
      headers: header,
    }

    fetch(`http://localhost:8080/movies`, requestoptions).then((response) => response.json()).then((data) => setmovies(data)).catch((error) => console.error(error))
  }, [])
  return (
    <>
      <h1>Movies</h1>
      <hr />

      <table className="table table-striped table-hover">
        <thead>
          <tr>
            <th>Movie</th>
            <th>Release-Date</th>
            <th>Rating</th>
          </tr>
        </thead>
        <tbody>
          {movies.map(m => (<tr key={m.id}>
            <td><Link to={`/movies/${m.id}`}>{m.title}</Link></td>
            <td>{m.release_date}</td>
            <td>{m.mpaa_rating}</td>
          </tr>))}
        </tbody>
      </table>
    </>
  )
}