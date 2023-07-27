package gensql

import (
	"database/sql"
	"math"
	"projects_template/internal/entity/global"

	"github.com/jmoiron/sqlx"
)

func Select[T any](tx *sqlx.Tx, sqlQuery string, params ...interface{}) ([]T, error) {
	data := make([]T, 0)

	err := tx.Select(&data, sqlQuery, params...)

	if err == nil && len(data) == 0 {
		err = sql.ErrNoRows
	}

	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, global.ErrNoData
	default:
		return nil, err
	}
}

func SelectNamed[T any](tx *sqlx.Tx, sqlQuery string, params map[string]interface{}) ([]T, error) {
	data := make([]T, 0)

	stmt, err := tx.PrepareNamed(sqlQuery)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Select(&data, params)
	if err != nil {
		return nil, err
	}

	if err == nil && len(data) == 0 {
		err = sql.ErrNoRows
	}

	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, global.ErrNoData
	default:
		return nil, err
	}
}

func SelectIn[T any](tx *sqlx.Tx, sqlQuery string, params interface{}) ([]T, error) {
	data := make([]T, 0)

	q, args, err := sqlx.In(sqlQuery, params)
	if err != nil {
		return nil, err
	}

	q = tx.Rebind(q)

	err = tx.Select(&data, q, args...)

	if len(data) == 0 && err == nil {
		err = sql.ErrNoRows
	}

	switch err {
	case nil:
		return data, nil
	case sql.ErrNoRows:
		return nil, global.ErrNoData
	default:
		return nil, err
	}
}

const SqlInBufferLimit = 1000

func SelectInOverLimit[T, K any](list []K, selectFunc func(resizedList []K) ([]T, error)) ([]T, error) {
	listLen := len(list)

	if listLen > SqlInBufferLimit {
		resultList := make([]T, 0, listLen)
		cycleCount := math.Ceil(float64(listLen) / float64(SqlInBufferLimit))

		start := 0
		end := SqlInBufferLimit

		for i := 0; i < int(cycleCount); i++ {
			tempResultList, err := selectFunc(list[start:end])
			switch err {
			case nil, global.ErrNoData:
				resultList = append(resultList, tempResultList...)

				start = end
				if end+SqlInBufferLimit > listLen {
					end = listLen
				} else {
					end += SqlInBufferLimit
				}
			default:
				return nil, err
			}
		}

		return resultList, nil
	} else {
		return selectFunc(list)
	}
}
