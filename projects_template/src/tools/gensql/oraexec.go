package gensql

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/mattn/go-oci8"
)

func OracleExecReturnID[T any](tx *sqlx.Tx, sqlInsert string, returnColumn, tableName string, insertParams ...interface{}) (id T, err error) {
	var res sql.Result

	res, err = tx.Exec(sqlInsert, insertParams...)
	if err != nil {
		return
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return
	}

	rowID := oci8.GetLastInsertId(lastInsertID)
	err = tx.QueryRow(fmt.Sprintf("select %s from %s where rowid = :1", returnColumn, tableName), rowID).Scan(&id)
	if err != nil {
		return
	}

	return
}

func OracleExecReturnIDNamed[T any](tx *sqlx.Tx, sqlInsert string, returnColumn, tableName string, insertParams map[string]interface{}) (id T, err error) {
	var res sql.Result

	stmt, err := tx.PrepareNamed(sqlInsert)
	if err != nil {
		return
	}
	defer stmt.Close()

	res, err = stmt.Exec(insertParams)
	if err != nil {
		return
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return
	}

	rowID := oci8.GetLastInsertId(lastInsertID)
	err = tx.QueryRow(fmt.Sprintf("select %s from %s where rowid = :1", returnColumn, tableName), rowID).Scan(&id)
	if err != nil {
		return
	}

	return
}
