package main

import (
	"log"
	"os"
	restapi "product_storage/external/restAPI"
	"product_storage/internal/transaction"
	"product_storage/rimport"
	"product_storage/tools/logger"
	"product_storage/tools/pgdb"
	"product_storage/uimport"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
}

var ()

func main() {
	version := os.Getenv("VERSION")
	db := pgdb.SqlxDB(os.Getenv("PG_URL"))

	log := logger.NewFileLogger("product_storage")
	log.Infoln("version", version)
	log.Debugln("pg", os.Getenv("PG_URL"))

	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	dbLog := logger.NewDBLog("product_storage", repo)

	useCase := uimport.NewUsecaseImports(log, dbLog, repo, repo.SessionManager)

	ginServer := restapi.NewGinServer(log, dbLog, useCase)
	ginServer.Run()

}
