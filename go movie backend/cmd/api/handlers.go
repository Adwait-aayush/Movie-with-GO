package main

import (
	"backend/internal/models"
	"errors"
	"log"
	"net/http"
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
	log.Printf("Received %s request to %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		app.ErrorJson(w, errors.New("only POST method is allowed"))
		return
	}

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
		app.ErrorJson(w, errors.New("password not found"))
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
	log.Println(tokens)
	app.WriteJson(w, http.StatusAccepted, tokens)
}

// func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
// 	for _, cookie := range r.Cookies() {
// 		if cookie.Name == app.auth.CookieName {
// 			claims := &Claims{}
// 			refreshToken := cookie.Value

//				_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
//					return []byte(app.JWTSecret), nil
//				})
//				if err != nil {
//					app.ErrorJson(w, errors.New("unauthorised"), http.StatusUnauthorized)
//					return
//				}
//				userID, err := strconv.Atoi(claims.Subject)
//				if err != nil {
//					app.ErrorJson(w, errors.New("unknown user"), http.StatusUnauthorized)
//					return
//				}
//				user, err := app.DB.GetUserByID(userID)
//				if err != nil {
//					app.ErrorJson(w, errors.New("unknown user"), http.StatusUnauthorized)
//					return
//				}
//				u := jwtUser{
//					ID:        user.ID,
//					FirstName: user.FirstName,
//					LastName:  user.LastName,
//				}
//				tokenPairs, err := app.auth.GenerateTokenPair(&u)
//				if err != nil {
//					app.ErrorJson(w, errors.New("ERROR IN TOKENS"), http.StatusUnauthorized)
//					return
//				}
//				http.SetCookie(w, app.auth.GetRefershCookie(tokenPairs.Refresh))
//				app.WriteJson(w, http.StatusOK, tokenPairs)
//			}
//		}
//	}
func (app *application) refreshToken(w http.ResponseWriter, r *http.Request) {
	var tokenPairs TokenPairs
	var foundCookie bool

	for _, cookie := range r.Cookies() {
		if cookie.Name == app.auth.CookieName {
			foundCookie = true
			claims := &Claims{}
			refreshToken := cookie.Value

			_, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
				return []byte(app.JWTSecret), nil
			})
			if err != nil {
				app.ErrorJson(w, errors.New("unauthorized"), http.StatusUnauthorized)
				return
			}

			if time.Now().Unix() > claims.ExpiresAt.Unix() {
				app.ErrorJson(w, errors.New("refresh token has expired"), http.StatusUnauthorized)
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

			tokenPairs, err = app.auth.GenerateTokenPair(&u)
			if err != nil {
				app.ErrorJson(w, errors.New("error generating tokens"), http.StatusInternalServerError)
				return
			}

			http.SetCookie(w, app.auth.GetRefershCookie(tokenPairs.Refresh))
			break
		}
	}

	if !foundCookie {
		app.ErrorJson(w, errors.New("unauthorized"), http.StatusUnauthorized)
		return
	}

	app.WriteJson(w, http.StatusOK, tokenPairs)
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
		Movie *models.Movie `json:"movie"`
		Genres []*models.Genre `json:"genres"`
	}{
		Movie: movie,
		Genres: genres,
	}
	_ = app.WriteJson(w, http.StatusOK, payload)

}
