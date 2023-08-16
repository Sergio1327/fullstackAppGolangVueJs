package gengin

import (
	"fmt"
	"product_storage/internal/entity/global"
	"product_storage/internal/entity/rpc"
	"product_storage/internal/transaction"
	"product_storage/tools/logger"

	"github.com/sirupsen/logrus"
)

type LoadDataFunc[T any] func(ts transaction.Session) (T, error)

func LoadData[T any](
	rc *rpc.Context,
	log *logrus.Logger,
	sessionManager transaction.SessionManager,
	f LoadDataFunc[T],
	stateName string,
	canErrNoData bool,
	needCommit bool,
) {
	ts := sessionManager.CreateSession()
	if err := ts.Start(); err != nil {
		log.Errorln(fmt.Sprintf("ошибка открытия транзакции; ошибка: %v", err))
		rc.ReturnError(global.ErrInternalError)
		return
	}
	defer ts.Rollback()

	data, err := f(ts)
	if err != nil {
		if err == global.ErrNoData && canErrNoData {
			rc.ReturnStateResult(stateName, nil)
			return
		}

		rc.GinContext.Error(logger.ErrLog(err))
		rc.ReturnError(err)
		return
	}

	if needCommit {
		err = ts.Commit()
		if err != nil {
			log.Errorln(fmt.Sprintf("ошибка при коммите; ошибка: %v", err))
			rc.ReturnError(global.ErrInternalError)
			return
		}
	}

	rc.ReturnStateResult(stateName, data)
}
