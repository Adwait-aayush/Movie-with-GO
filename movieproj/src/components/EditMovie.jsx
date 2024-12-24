import { useEffect, useState } from "react"
import { useNavigate, useOutletContext, useParams } from "react-router-dom"
import "./Editmovie.css"

export default function EditMovie() {
    const { jwttok } = useOutletContext()

    const [movie, setMovie] = useState({
        id: 0,
        title: "",
        release_date: "",
        runtime: "",
        mpaa_rating: "",
        description: "",
    })
    let { id } = useParams()
    const navigate = useNavigate()

    useEffect(() => {
        if (jwttok === "") {
            navigate("/login")
            return
        }
    }, [jwttok, navigate])

    const handleSubmit = (e) => {
        e.preventDefault()
    }

    const handleChange = (e) => {
        let name = e.target.name
        let value = e.target.value
        setMovie({
            ...movie,
            [name]: value,
        })
    }

    return (
        <div className="edit-movie-container">
            <h1>ADD/EDIT Movie</h1>
            <hr />
            <pre>{JSON.stringify(movie,null,3)}</pre>
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
                <button type="submit" className="submit-btn">Save Changes</button>
            </form>
        </div>
    )
}
