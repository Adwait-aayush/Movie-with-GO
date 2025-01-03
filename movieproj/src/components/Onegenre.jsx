import { useEffect, useState } from "react";
import { Link, useLocation, useParams } from "react-router-dom";

export default function OneGenre() {
    const location = useLocation();
    const { genreName } = location.state;

    const [movies, setMovies] = useState([]);
    let { id } = useParams();

    useEffect(() => {
        const headers = new Headers();
        headers.append("Content-Type", "application/json");

        const requestOptions = {
            method: "GET",
            headers: headers,
        };

        fetch(`http://localhost:8080/movies/genres/${id}`, requestOptions)
            .then((response) => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error);
                } else {
                    setMovies(data);
                }
            })
            .catch((err) => {
                console.log(err);
            });
    }, [id]);

    return (
        <>
            <h2>Genre: {genreName}</h2>
            <hr />
            {movies.length > 0 ? (
                <table className="table table-striped table-hover">
                    <thead>
                        <tr>
                            <th>Movie</th>
                            <th>Release Date</th>
                            <th>Rating</th>
                        </tr>
                    </thead>
                    <tbody>
                        {movies.map((m) => {
                            return (
                                <tr key={m.id}>
                                    <td>
                                        <Link to={`/movies/${m.id}`}>{m.title}</Link>
                                    </td>
                                    <td>{m.release_date}</td>
                                    <td>{m.mpaa_rating}</td>
                                </tr>
                            );
                        })}
                    </tbody>
                </table>
            ) : (
                <div className="alert alert-info">NO Movies yet</div>
            )}
        </>
    );
}
