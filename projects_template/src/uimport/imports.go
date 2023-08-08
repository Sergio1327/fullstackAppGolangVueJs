package uimport

import (
	"os"
	"product_storage/config"
	"product_storage/internal/transaction"
	"product_storage/internal/usecase"
	"product_storage/rimport"

	"github.com/sirupsen/logrus"
)

type UsecaseImports struct {
	Config         config.Config
	SessionManager transaction.SessionManager
	Usecase        Usecase
}

func NewUsecaseImports(
	log *logrus.Logger,
	dblog *logrus.Logger,
	ri rimport.RepositoryImports,
	sessionManager transaction.SessionManager,
) UsecaseImports {
	config, err := config.NewConfig(os.Getenv("CONF_PATH"))
	if err != nil {
		log.Fatalln(err)
	}

	ui := UsecaseImports{
		Config:         config,
		SessionManager: sessionManager,

		Usecase: Usecase{
			Logger:  usecase.NewLogger(log, ri),
			Product: usecase.NewProduct(log, dblog, ri),
		},
	}

	return ui
}
