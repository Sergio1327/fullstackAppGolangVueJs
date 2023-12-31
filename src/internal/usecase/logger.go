package usecase

import (
	"github.com/sirupsen/logrus"
	"product_storage/internal/entity/log"
	"product_storage/rimport"
)

type Logger struct {
	log *logrus.Logger
	ri  rimport.RepositoryImports
}

func NewLogger(log *logrus.Logger,
	ri rimport.RepositoryImports,
) *Logger {
	return &Logger{
		log: log,
		ri:  ri,
	}
}

func (u *Logger) logPrefix() string {
	return "[logger_usecase]"
}

func (u *Logger) SpecialFields() []string {
	return []string{"product_id", "variant_id", "stock_id"}
}

func (u *Logger) SaveLog(row log.Row) error {
	ts := u.ri.SessionManager.CreateSession()
	if err := ts.Start(); err != nil {
		u.log.Errorln(u.logPrefix(), "не удается стартовать транзакцию", err)
		return err
	}
	defer ts.Rollback()

	if row.Details != nil && len(row.Details) > 0 {
		logID, err := u.ri.Repository.Logger.SaveLogWithReturnID(ts, row)
		if err != nil {
			u.log.Errorln(u.logPrefix(), "не удается сохранить данные в лог", err)
			return err
		}

		if err := u.ri.Repository.Logger.SaveLogDetails(ts, logID, row.Details); err != nil {
			u.log.Errorln(u.logPrefix(), "не удается сохранить данные в лог", err)
			return err
		}

	} else {
		if err := u.ri.Repository.Logger.SaveLog(ts, row); err != nil {
			u.log.Errorln(u.logPrefix(), "не удается сохранить данные в детали лога", err)
			return err
		}
	}

	if err := ts.Commit(); err != nil {
		u.log.Errorln(u.logPrefix(), "не удается зафиксировать транзакцию", err)
		return err
	}

	return nil
}
