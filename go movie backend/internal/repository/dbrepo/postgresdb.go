package dbrepo

import (
	"backend/internal/models"
	"context"
	"database/sql"
	"time"
)

type Postgresdbrepo struct {
	Db *sql.DB
}

const dbtimeout = time.Second * 3

func (m *Postgresdbrepo) Connection() *sql.DB {
	return m.Db
}

func (m *Postgresdbrepo) AllMovies() ([]*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbtimeout)
	defer cancel()

	query := `
    select id, title, release_date, runtime, mpaa_rating, description, coalesce(image, ''),
    created_at, updated_at
    from movies
    order by title
    `

	rows, err := m.Db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*models.Movie
	for rows.Next() {
		var movie models.Movie
		err = rows.Scan(
			&movie.ID,
			&movie.Title,
			&movie.ReleaseDate,
			&movie.Runtime,
			&movie.Rating,
			&movie.Description,
			&movie.Image,
			&movie.CreatedAt,
			&movie.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return movies, nil

}

func (m *Postgresdbrepo) OneMovie(id int) (*models.Movie, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbtimeout)
	defer cancel()

	query := ` select id, title, release_date, runtime, mpaa_rating, description, coalesce(image, ''),
    created_at, updated_at
    from movies
    where id=$1
    `
	row := m.Db.QueryRowContext(ctx, query, id)
	var movie models.Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	query = `select g.id,g.genre from movies_genres mg
	left join genres g on (mg.genre_id=g.id)
	where mg.movie_id=$1
	 order by g.genre`
	rows, err := m.Db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var genres []*models.Genre
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, err
		}
		genres = append(genres, &g)
	}
	movie.Genre = genres
	return &movie, nil

}

func (m *Postgresdbrepo) OneMovieForEdit(id int) (*models.Movie,[]*models.Genre ,error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbtimeout)
	defer cancel()

	query := ` select id, title, release_date, runtime, mpaa_rating, description, coalesce(image, ''),
    created_at, updated_at
    from movies
    where id=$1
    `
	row := m.Db.QueryRowContext(ctx, query, id)
	var movie models.Movie
	err := row.Scan(
		&movie.ID,
		&movie.Title,
		&movie.ReleaseDate,
		&movie.Runtime,
		&movie.Rating,
		&movie.Description,
		&movie.Image,
		&movie.CreatedAt,
		&movie.UpdatedAt,
	)
	if err != nil {
		return nil, nil,err
	}

	query = `select g.id,g.genre from movies_genres mg
	left join genres g on (mg.genre_id=g.id)
	where mg.movie_id=$1
	 order by g.genre`
	rows, err := m.Db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, nil,err
	}
	defer rows.Close()
	var genres []*models.Genre
	var genresArray []int
	for rows.Next() {
		var g models.Genre
		err := rows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil,err
		}
		genres = append(genres, &g)
		genresArray = append(genresArray, g.ID)
	}
	movie.Genre = genres
	movie.GenresArray=genresArray
	var allGenres []*models.Genre

	query=`select id,genre from genres order by genre`
	grows,err:=m.Db.QueryContext(ctx,query)
	if err != nil {
		return nil, nil,err
	}
	defer grows.Close()
	for grows.Next() {
		var g models.Genre
		err := grows.Scan(
			&g.ID,
			&g.Genre,
		)
		if err != nil {
			return nil, nil,err
		}
		allGenres=append(allGenres,&g )
	}
	return &movie,allGenres,nil

}

func (m *Postgresdbrepo) GetUserByEmail(email string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbtimeout)
	defer cancel()

	query := `select id,email,first_name,last_name,password,created_at, updated_at from users where email=$1`
	var user models.User
	row := m.Db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (m *Postgresdbrepo) GetUserByID(id int) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbtimeout)
	defer cancel()

	query := `select id,email,first_name,last_name,password,created_at, updated_at from users where id=$1`
	var user models.User
	row := m.Db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
