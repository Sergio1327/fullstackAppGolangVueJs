package main

import (
	"backend/db"
	"backend/internal/server"
	"log"
)

func main() {
	db, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer db.DB.Close()

	err = server.Run(db)
	if err != nil {
		log.Fatal(err)
	}
}
