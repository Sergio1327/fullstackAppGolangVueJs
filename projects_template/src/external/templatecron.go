package external

import (
	"projects_template/uimport"

	"github.com/sirupsen/logrus"
)

type Cron struct {
	// вначале системные объекты - логи, конфиги
	log   *logrus.Logger
	dblog *logrus.Logger
	// далее usecase
	uimport.UsecaseImports
}

func NewCron(log *logrus.Logger,
	dblog *logrus.Logger,
	// далее usecase
	u uimport.UsecaseImports) *Cron {
	return &Cron{
		log:            log,
		dblog:          dblog,
		UsecaseImports: u,
	}
}

func (c *Cron) Run() {

}
