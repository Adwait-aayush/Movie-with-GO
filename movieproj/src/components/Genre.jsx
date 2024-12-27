import { useEffect } from "react";
import { useState } from "react";
import { Link } from "react-router-dom";

export default function Genre() {
    const [genres, setGenres] = useState([]);

    useEffect(() => {
        const headers = new Headers();
        headers.append('Content-Type', 'application/json');
        const reqOptions = {
            method: 'GET',
            headers: headers,
        };
        fetch(`http://localhost:8080/genres`, reqOptions)
            .then(response => response.json())
            .then((data) => {
                if (data.error) {
                    console.log(data.error);
                } else {
                    setGenres(data);
                }
            })
            .catch(err => {
                console.log(err);
            });
    }, []);

    return (
        <>
            <h1>Welcome to Genre Page</h1>
            <hr />
            <div className="list-group">
                {genres.map((g) => {
                    return (
                        <Link
                            key={g.id}
                            className="list-group-item list-group-item-action"
                            to={`/genre/${g.id}`}
                            state={{
                                genreName: g.genre,
                            }}
                        >
                            {g.genre}
                        </Link>
                    );
                })}
            </div>
        </>
    );
}
