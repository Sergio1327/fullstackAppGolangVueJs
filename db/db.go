package db

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"
)

type Db struct {
	DB *sqlx.DB
}

func Connect() (*Db, error) {
	db, err := sqlx.Open("postgres", "host=127.0.0.1 port=5432 user=test_db password=test_db dbname=test_db sslmode=disable")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if err := db.Ping(); err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("succesfully connected to db")
	return &Db{
		DB: db,
	}, nil
}

func (d *Db) Close() error {
	return d.DB.Close()
}

func (d *Db) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return d.DB.Query(query, args...)
}

func (d *Db) QueryRow(query string, args ...interface{}) *sql.Row {
	return d.DB.QueryRow(query, args...)
}

func (d *Db) Exec(query string, args ...interface{}) (sql.Result, error) {
	return d.DB.Exec(query, args...)
}
