import { useParams } from "react-router-dom"
import { useState } from "react"
import { useEffect } from "react"

export default function Movie(){
    const [Movie,setMovie]=useState({})
    const {id}=useParams()
    useEffect(()=>{
        const headers= new Headers()
        headers.append('Content-Type','application/json')
        const requestoptions={
            method:'GET',
            headers:headers
        }
        fetch(`http://localhost:8080/movies/${id}`,requestoptions)
        .then(response=>response.json())
        .then(data=>{setMovie(data)})
        .catch(e=>{
            console.log(e)
        })
    },[id])
if (Movie.genres){
    Movie.genres=Object.values(Movie.genres)
}
else{
    Movie.genres=[]
}

    return(
        <>
        <h1>Movie:{Movie.title}</h1>
        <hr /> 
        {Movie.genres.map((g)=>(
            <span key={g.genre} className="badge bg-secondary me-2">{g.genre}</span>
        ))}
        <p>Release Date: {Movie.release_date}</p>
        <p>Run Time: {Movie.runtime}</p>
        <hr />
        {Movie.image !=="" &&
        <div className="mb-3">
            <img src={`https://image.tmdb.org/t/p/w200/${Movie.image}`} alt="Poster" />
            </div>}
        <p>{Movie.description}</p>

        </>
    )
}