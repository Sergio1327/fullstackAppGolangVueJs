package main

import (
	"os"

	restapi "product_storage/external/restAPI"
	"product_storage/internal/transaction"
	"product_storage/rimport"
	"product_storage/tools/logger"
	"product_storage/tools/pgdb"
	"product_storage/uimport"
)

func init() {

}

func main() {
	version := os.Getenv("VERSION")
	pgURL := os.Getenv("PG_URL")

	db := pgdb.SqlxDB(pgURL)
	defer db.Close()

	log := logger.NewFileLogger("product_storage")
	log.Infoln("version", version)
	log.Info("pg", os.Getenv("PG_URL"))

	sm := transaction.NewSQLSessionManager(db)
	repo := rimport.NewRepositoryImports(sm)

	dbLog := logger.NewDBLog("product_storage", repo)

	useCase := uimport.NewUsecaseImports(log, dbLog, repo, repo.SessionManager)

	ginServer := restapi.NewGinServer(log, dbLog, useCase)
	ginServer.Run()
}
