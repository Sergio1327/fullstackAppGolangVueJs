package main

import (
	"os"
	"product_storage/tools/logger"
)

var (
	version string = os.Getenv("VERSION")
)

func main() {
	log := logger.NewNoFileLogger("product_storage")
	log.Infoln("version", version)
	log.Debugln("pg", os.Getenv("PG_URL"))
}
