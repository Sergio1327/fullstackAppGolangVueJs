package usecase

import (
	"product_storage/internal/entity/log"
	"product_storage/rimport"
	"product_storage/tools/sqlnull"

	"github.com/sirupsen/logrus"
)

type LoggerUsecase struct {
	log *logrus.Logger
	ri  rimport.RepositoryImports
}

func NewLogger(log *logrus.Logger,
	ri rimport.RepositoryImports,
) *LoggerUsecase {
	return &LoggerUsecase{
		log: log,
		ri:  ri,
	}
}

func (u *LoggerUsecase) logPrefix() string {
	return "[logger_usecase]"
}

func (u *LoggerUsecase) SpecialFields() []string {
	return []string{"c_id", "se_id", "oper_login"}
}

func (u *LoggerUsecase) SaveLog(row log.Row) error {
	ts := u.ri.SessionManager.CreateSession()
	if err := ts.Start(); err != nil {
		u.log.Errorln(u.logPrefix(), "не удается стартовать транзакцию", err)
		return err
	}
	defer ts.Rollback()

	var (
		contractID sqlnull.NullInt64
		seID       sqlnull.NullInt64
		operLogin  sqlnull.NullString
	)

	if data, exists := row.SpecialFields["c_id"]; exists {
		if data.Type == "int" {
			contractID.Scan(data.Value)
		}
	}

	if data, exists := row.SpecialFields["se_id"]; exists {
		if data.Type == "int" {
			seID.Scan(data.Value)
		}
	}

	if data, exists := row.SpecialFields["oper_login"]; exists {
		if data.Type == "string" {
			operLogin.Scan(data.Value)
		}
	}

	if row.Details != nil && len(row.Details) > 0 {
		logID, err := u.ri.Repository.Logger.SaveLogWithReturnID(ts, row, contractID, seID, operLogin)
		if err != nil {
			u.log.Errorln(u.logPrefix(), "не удается сохранить данные в лог", err)
			return err
		}

		if err := u.ri.Repository.Logger.SaveLogDetails(ts, logID, row.Details); err != nil {
			u.log.Errorln(u.logPrefix(), "не удается сохранить данные в лог", err)
			return err
		}

	} else {
		if err := u.ri.Repository.Logger.SaveLog(ts, row, operLogin); err != nil {
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
