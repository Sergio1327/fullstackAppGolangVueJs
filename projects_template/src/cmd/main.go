package main

import "projects_template/tools/logger"

var (
	version string
)

func main() {
	log := logger.NewFileLogger("projects_template")
	log.Infoln("version", version)
}
