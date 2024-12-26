
import { useEffect, useState } from "react"
import { Link, useNavigate, useOutletContext } from "react-router-dom";

export default function Manage(){
     const [movies, setmovies] = useState([])

     const {jwttok}=useOutletContext();
     const navigate=useNavigate();

    
      useEffect(() => {
        if (jwttok==="") {
       navigate("/login");
        }
        const header = new Headers();
        header.append('Content-Type', 'application/json');
        header.append('Authorization', `Bearer ${jwttok}`);
        const requestoptions = {
          method: 'GET',
          headers: header,
        }
    
        fetch(`http://localhost:8080/admin/movies`, requestoptions).then((response) => response.json()).then((data) => setmovies(data)).catch((error) => console.error(error))
      }, [jwttok,navigate])
    return(
        <>
        <h1>
            Manage-Catalogue
        </h1>
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
            <td><Link to={`/admin/movie/${m.id}`}>{m.title}</Link></td>
            <td>{m.release_date}</td>
            <td>{m.mpaa_rating}</td>
          </tr>))}
        </tbody>
      </table>
        </>
    )
}