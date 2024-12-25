package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func (app *application) router() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)
	mux.Get("/", app.Home)
	mux.Get("/movies", app.Movies)
	mux.Post("/authenticate", app.authenticate)
	mux.Get("/movies/{id}", app.GetMovie)
	mux.Get("/refresh", app.refreshToken)
	mux.Get("/logout", app.logout)
	return mux
}
