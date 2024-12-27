package main

import (
	"backend/internal/graph"
	"backend/internal/models"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/golang-jwt/jwt/v4"
)

func (app *application) Home(w http.ResponseWriter, r *http.Request) {
	var payload = struct {
		Status  string `json:"Status"`
		Message string `json:"Message"`
		Version string `json:"Version"`
	}{
		Status:  "OK",
		Message: "Hello, World!",
		Version: "1.0",
	}
	app.WriteJson(w, http.StatusOK, payload)
}
func (app *application) Movies(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	app.WriteJson(w, http.StatusOK, movies)

}

func (app *application) authenticate(w http.ResponseWriter, r *http.Request) {

	//read data
	var requestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	err := app.ReadJson(w, r, &requestPayload)
	if err != nil {
		app.ErrorJson(w, errors.New("data not read properly"), http.StatusUnauthorized)
		return
	}

	//validate
	user, err := app.DB.GetUserByEmail(requestPayload.Email)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	//check pass
	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		app.ErrorJson(w, errors.New("password not found"), http.StatusBadRequest)
		return
	}

	//create user
	u := jwtUser{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
	}
	//generate token
	tokens, err := app.auth.GenerateTokenPair(&u)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	refreshCookie := app.auth.GetRefershCookie(tokens.Refresh)
	http.SetCookie(w, refreshCookie)

	app.WriteJson(w, http.StatusAccepted, tokens)
}

func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			claims := &Claims{}
			refreshToken := cookie.Value

			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.ErrorJson(w, errors.New("unauthorised"), http.StatusUnauthorized)
				return
			}
			userID, err := strconv.Atoi(claims.Subject)
			if err != nil {
				app.ErrorJson(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}
			user, err := app.DB.GetUserByID(userID)
			if err != nil {
				app.ErrorJson(w, errors.New("unknown user"), http.StatusUnauthorized)
				return
			}
			u := jwtUser{
				ID:        user.ID,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			}
			tokenPairs, err := app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.ErrorJson(w, errors.New("ERROR IN TOKENS"), http.StatusUnauthorized)
				return
			}
			http.SetCookie(w, app.auth.GetRefershCookie(tokenPairs.Refresh))
			app.WriteJson(w, http.StatusOK, tokenPairs)

		}
	}
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, app.auth.GetExpiredRefershCookie())
	w.WriteHeader(http.StatusAccepted)
}

func (app *application) MovieCatalogue(w http.ResponseWriter, r *http.Request) {
	movies, err := app.DB.AllMovies()
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	app.WriteJson(w, http.StatusOK, movies)

}

func (app *application) GetMovie(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.ErrorJson(w, errors.New("invalid movie id"), http.StatusBadRequest)
	}
	movie, err := app.DB.OneMovie(movieID)
	if err != nil {
		app.ErrorJson(w, errors.New("movie not found"), http.StatusNotFound)
	}
	_ = app.WriteJson(w, http.StatusOK, movie)
}

func (app *application) MovieForEdit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	movieID, err := strconv.Atoi(id)
	if err != nil {
		app.ErrorJson(w, errors.New("invalid movie id"), http.StatusBadRequest)
		return
	}
	movie, genres, err := app.DB.OneMovieForEdit(movieID)
	if err != nil {
		app.ErrorJson(w, errors.New("movie not found"), http.StatusNotFound)
		return
	}
	var payload = struct {
		Movie  *models.Movie   `json:"movie"`
		Genres []*models.Genre `json:"genres"`
	}{
		Movie:  movie,
		Genres: genres,
	}
	_ = app.WriteJson(w, http.StatusOK, payload)

}
func (app *application) AllGenres(w http.ResponseWriter, r *http.Request) {
	genres, err := app.DB.AllGenres()
	if err != nil {
		app.ErrorJson(w, err)
	}
	_ = app.WriteJson(w, http.StatusOK, genres)
}

func (app *application) InsertMovie(w http.ResponseWriter, r *http.Request) {
	var movie models.Movie
	err := app.ReadJson(w, r, &movie)
	if err != nil {
		app.ErrorJson(w, err)
		return

	}
	//image

	movie = app.GetPoster(movie)
	movie.CreatedAt = time.Now()
	movie.UpdatedAt = time.Now()

	newID, err := app.DB.InsertMovie(movie)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	//genre
	err = app.DB.UpdateMovieGenres(newID, movie.GenresArray)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Movie inserted successfully",
	}
	app.WriteJson(w, http.StatusAccepted, resp)
}

func (app *application) GetPoster(movie models.Movie) models.Movie {
	type TheMovieDB struct {
		Page    int `json:"page"`
		Results []struct {
			PosterPath string `json:"poster_path"`
		} `json:"results"`
		TotalPages int `json:"total_pages"`
	}

	client := &http.Client{}
	theUrl := fmt.Sprintf("http://api.themoviedb.org/3/search/movie?api_key=%s", app.APIKey)
	req, err := http.NewRequest("GET", theUrl+"&query="+url.QueryEscape(movie.Title), nil)
	if err != nil {
		log.Println(err)
		return movie
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return movie
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return movie
	}
	var responseobject TheMovieDB
	json.Unmarshal(bodyBytes, &responseobject)
	if len(responseobject.Results) > 0 {
		movie.Image = responseobject.Results[0].PosterPath
	}
	return movie

}

func (app *application) Updatemovie(w http.ResponseWriter, r *http.Request) {
	var payload models.Movie
	err := app.ReadJson(w, r, &payload)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	movie, err := app.DB.OneMovie(payload.ID)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	movie.Title = payload.Title
	movie.ReleaseDate = payload.ReleaseDate
	movie.Description = payload.Description
	movie.Rating = payload.Rating
	movie.Runtime = payload.Runtime
	movie.UpdatedAt = time.Now()

	err = app.DB.UpdateMovie(*movie)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}

	resp := JSONResponse{
		Error:   false,
		Message: "Movie updated successfully",
	}
	app.WriteJson(w, http.StatusAccepted, resp)
}
func (app *application) DeleteMovie(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	err = app.DB.DeleteMovie(id)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	resp := JSONResponse{
		Error:   false,
		Message: "Movie deleted successfully",
	}
	app.WriteJson(w, http.StatusAccepted, resp)
}

func (app *application) AllmoviesBygenre(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	movies, err := app.DB.AllMovies(id)
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	app.WriteJson(w, http.StatusOK, movies)
}

func (app *application) graphqlmovs(w http.ResponseWriter, r *http.Request) {
	movies, _ := app.DB.AllMovies()

	q, _ := io.ReadAll(r.Body)
	query := string(q)

	g := graph.New(movies)

	g.QueryString = query
	resp, err := g.Query()
	if err != nil {
		app.ErrorJson(w, err)
		return
	}
	j, _ := json.MarshalIndent(resp, " ", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(j)

}
