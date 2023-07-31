package rimport

import (
	"log"
	"os"
	"projects_template/config"
	"projects_template/internal/repository/postgresql"
	"projects_template/internal/transaction"
)

type RepositoryImports struct {
	Config         config.Config
	SessionManager transaction.SessionManager
	Repository     Repository
}

func NewRepositoryImports(
	sessionManager transaction.SessionManager,
) RepositoryImports {
	config, err := config.NewConfig(os.Getenv("CONF_PATH"))
	if err != nil {
		log.Fatalln(err)
	}

	return RepositoryImports{
		Config:         config,
		SessionManager: sessionManager,
		Repository: Repository{
			Product: postgresql.NewPostgresProductRepository(),
		},
	}
}
