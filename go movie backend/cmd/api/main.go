package main

import (
	"backend/internal/repository"
	"backend/internal/repository/dbrepo"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

type application struct {
	DNS          string
	Domain       string
	auth         Auth
	JWTSecret    string
	JWTIssuer    string
	JWTAudience  string
	CookieDomain string

	DB repository.Databaserepo
	APIKey string
}

func main() {
	var app application

	flag.StringVar(&app.DNS, "dns", "host=localhost port=5432 user=postgres password=postgres dbname=movies sslmode=disable timezone=UTC connect_timeout=5", "postgress datab connection")
	flag.StringVar(&app.JWTSecret, "jwt-secret", "verysecret", "Signinsecret")
	flag.StringVar(&app.JWTIssuer, "jwt-issuer", "example.com", "Signinissuer")
	flag.StringVar(&app.JWTAudience, "jwt-audience", "example.com", "Signinaudience")
	flag.StringVar(&app.CookieDomain, "cookie-domain", "localhost", "cookie domain")
	flag.StringVar(&app.Domain, "domain", "example.com", "domain")
	flag.StringVar(&app.APIKey, "api-key", "0b124068ce1417120e20cb21c5e5213e", "api key")
	
	flag.Parse()

	conn, err := app.connectToDB()
	if err != nil {
		log.Fatal(err)
	}
	app.DB = &dbrepo.Postgresdbrepo{Db: conn}
	defer app.DB.Connection().Close()
	app.auth=Auth{
		Issuer: app.JWTIssuer,
		Audience: app.JWTAudience,
		Secret: app.JWTSecret,
		TokenExpiry: time.Minute*15,
		RefreshExpiry: time.Hour*24,
		CookiePath: "/",
		CookieName: "refresh_token",
		CookieDomain: app.CookieDomain,

	}

	fmt.Println("listening on port 8080")
	err = http.ListenAndServe(":8080", app.router())
	if err != nil {
		log.Fatal(err)
	}
}
