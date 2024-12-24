package main

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

func OpenDb(dsn string) (*sql.DB, error){
	db, err := sql.Open("pgx", dsn)
	if err!=nil {
		return nil,err
	}
	err=db.Ping()
	if err!=nil {
		return nil,err
	}
	return db, nil
}
func (app *application) connectToDB() (*sql.DB,error){
	connection,err:=OpenDb(app.DNS)
	if err!=nil {
		return nil,err
	}
	log.Println("Connected to Database")
	return connection,nil
}
