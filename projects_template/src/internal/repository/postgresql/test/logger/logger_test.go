package logger_test

import (
	"product_storage/internal/entity/log"
	"product_storage/internal/repository/postgresql"
	"product_storage/internal/transaction"
	"product_storage/tools/pgdb"
	"product_storage/tools/sqlnull"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSaveLog(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()

	repo := postgresql.NewLoggerRepository()

	ts := transaction.NewSQLSession(db)
	err := ts.Start()
	r.NoError(err)
	defer ts.Rollback()

	log := log.Row{
		Time:    time.Now(),
		Flag:    "INFO",
		Message: "fkjdskjs",
		Module:  "tests",
		File:    sqlnull.NewString("logger_test.go"),
		Line:    sqlnull.NewString("36"),
	}
	err = repo.SaveLog(ts, log)
	r.NoError(err)
}

func TestSaveLogWithReturnID(t *testing.T) {
	r := require.New(t)

	db := pgdb.SqlxDB("dbname=test_db user=test_db password=test_db host=127.0.0.1 port=5432 sslmode=disable")
	defer db.Close()

	repo := postgresql.NewLoggerRepository()

	ts := transaction.NewSQLSession(db)
	err := ts.Start()
	r.NoError(err)
	defer ts.Rollback()

	log := log.Row{
		Time:    time.Now(),
		Flag:    "INFO",
		Message: "fkjdskjs",
		Module:  "tests",
		File:    sqlnull.NewString("logger_test.go"),
		Line:    sqlnull.NewString("36"),
	}

	logID, err := repo.SaveLogWithReturnID(ts, log)
	r.NoError(err)
	r.NotZero(logID)
}
