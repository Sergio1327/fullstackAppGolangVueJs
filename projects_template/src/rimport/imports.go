package rimport

import (
	"log"
	"os"
	"product_storage/config"
	"product_storage/internal/repository/postgresql"
	"product_storage/internal/transaction"
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
			Product: postgresql.NewProduct(),
		},
	}
}
