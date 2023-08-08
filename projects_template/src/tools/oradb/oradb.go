package oradb

import (
	"log"

	"github.com/jmoiron/sqlx"
)

// AlterOracleFloatParseFix fix for float
const AlterOracleFloatParseFix = `ALTER SESSION SET NLS_NUMERIC_CHARACTERS = '. '`

// NewOracleDB get db link
func NewOracleDB(oracleURL string) *sqlx.DB {
	db, err := sqlx.Connect("oci8", oracleURL)
	if err != nil {
		log.Fatalln(err)
	}
	return db
}
