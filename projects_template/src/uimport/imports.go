package uimport

import (
	"os"
	"projects_template/bimport"
	"projects_template/config"
	"projects_template/internal/transaction"
	"projects_template/internal/usecase"
	"projects_template/internal/usecase/product"
	"projects_template/rimport"

	"github.com/sirupsen/logrus"
)

type UsecaseImports struct {
	Config         config.Config
	SessionManager transaction.SessionManager
	Usecase        Usecase
	*bimport.BridgeImports
}

func NewUsecaseImports(
	log *logrus.Logger,
	dblog *logrus.Logger,
	ri rimport.RepositoryImports,
	bi *bimport.BridgeImports,
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
			Logger:   usecase.NewLogger(log, ri),
			ProdcutUsecase: product.NewProductUseCaseImpl(log,dblog,ri),
		},
		BridgeImports: bi,
	}

	return ui
}
