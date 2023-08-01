package logger_test

import (
	"os"
	"product_storage/config"
	"product_storage/internal/entity/log"
	"product_storage/internal/repository/oracle"
	"product_storage/internal/transaction"
	"product_storage/tools/oradb"
	"product_storage/tools/sqlnull"

	"testing"

	"github.com/bxcodec/faker"
	"github.com/stretchr/testify/require"
)

func TestSaveLog(t *testing.T) {
	r := require.New(t)
	conf, err := config.NewConfig(os.Getenv("CONF_PATH"))
	r.NoError(err)

	db := oradb.NewOracleDB(conf.OracleConnectString())
	r.NoError(db.Ping())

	repo := oracle.NewLogger()

	ts := transaction.NewSQLSession(db)
	err = ts.Start()
	r.NoError(err)
	defer ts.Rollback()

	var (
		row        log.Row
		contractID sqlnull.NullInt64
		seID       sqlnull.NullInt64
		operLogin  sqlnull.NullString
	)

	faker.FakeData(&row)
	faker.FakeData(&contractID)
	faker.FakeData(&seID)
	faker.FakeData(&operLogin)

	row.Flag = "INFO"
	row.Line.Scan("10")

	err = repo.SaveLog(ts, row, contractID, seID, operLogin)
	r.NoError(err)
}

func TestSaveLogWithReturnID(t *testing.T) {
	r := require.New(t)
	conf, err := config.NewConfig(os.Getenv("CONF_PATH"))
	r.NoError(err)

	db := oradb.NewOracleDB(conf.OracleConnectString())
	r.NoError(db.Ping())

	repo := oracle.NewLogger()

	ts := transaction.NewSQLSession(db)
	err = ts.Start()
	r.NoError(err)
	defer ts.Rollback()

	var (
		row        log.Row
		contractID sqlnull.NullInt64
		seID       sqlnull.NullInt64
		operLogin  sqlnull.NullString
	)

	faker.FakeData(&row)
	faker.FakeData(&contractID)
	faker.FakeData(&seID)
	faker.FakeData(&operLogin)

	row.Flag = "INFO"
	row.Line.Scan("10")

	id, err := repo.SaveLogWithReturnID(ts, row, contractID, seID, operLogin)
	r.NoError(err)
	r.NotEmpty(id)
}

func TestSaveLogDetails(t *testing.T) {
	r := require.New(t)
	conf, err := config.NewConfig(os.Getenv("CONF_PATH"))
	r.NoError(err)

	db := oradb.NewOracleDB(conf.OracleConnectString())
	r.NoError(db.Ping())

	repo := oracle.NewLogger()

	ts := transaction.NewSQLSession(db)
	err = ts.Start()
	r.NoError(err)
	defer ts.Rollback()

	var (
		row        log.Row
		contractID sqlnull.NullInt64
		seID       sqlnull.NullInt64
		operLogin  sqlnull.NullString
	)

	faker.FakeData(&row)
	faker.FakeData(&contractID)
	faker.FakeData(&seID)
	faker.FakeData(&operLogin)

	row.Flag = "INFO"
	row.Line.Scan("10")

	id, err := repo.SaveLogWithReturnID(ts, row, contractID, seID, operLogin)
	r.NoError(err)
	r.NotEmpty(id)

	r.NoError(repo.SaveLogDetails(ts, id, row.Details))
}
