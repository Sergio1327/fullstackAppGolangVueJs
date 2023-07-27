package genproxy

import (
	"projects_template/internal/entity/global"

	"github.com/sirupsen/logrus"
)

type LoadDataFunc[T any] func() (T, error)

func LoadData[T any](f LoadDataFunc[T], l *logrus.Logger, errMsg ...interface{}) (T, error) {
	data, err := f()
	switch err {
	case nil:
		return data, nil
	case global.ErrNoData:
		return data, global.ErrNoData
	default:
		errMsg = append(errMsg, []interface{}{"ошибка:", err}...)
		l.Errorln(errMsg...)
		return data, global.ErrInternalError
	}
}
