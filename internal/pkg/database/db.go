package database

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

func NewPostgreSQLdb(str string) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", str)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	log.Println("Succesfully connected to database")
	return db, nil
}
