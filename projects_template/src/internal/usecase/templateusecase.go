package usecase

import (
	"fmt"
	"projects_template/bimport"
	"projects_template/internal/entity/global"
	"projects_template/internal/entity/template"
	"projects_template/internal/transaction"
	"projects_template/rimport"

	"github.com/sirupsen/logrus"
)

type TemplateUsecase struct {
	// вначале системные объекты - логи, конфиги
	log   *logrus.Logger
	dblog *logrus.Logger
	// далее репозитории
	rimport.RepositoryImports
	*bimport.BridgeImports
}

func NewTemplate(
	log, dblog *logrus.Logger,
	ri rimport.RepositoryImports,
	bi *bimport.BridgeImports,
) *TemplateUsecase {
	return &TemplateUsecase{
		log:               log,
		dblog:             dblog,
		RepositoryImports: ri,
		BridgeImports:     bi,
	}
}

func (u *TemplateUsecase) AwesomePublicMethod(ts transaction.Session, id int) (data template.TemplateObject, err error) {
	lf := logrus.Fields{"id": id}

	data, err = u.Repository.Template.FindTemplateObj(ts, id)
	if err != nil {
		u.dblog.WithFields(lf).Errorln(fmt.Sprintf("не удалось получить объект шаблона; ошибка: %v", err))
		err = global.ErrInternalError
		return
	}

	return
}
