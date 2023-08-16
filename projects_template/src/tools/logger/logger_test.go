package logger

import (
	"os"
	"product_storage/config"
	"product_storage/internal/transaction"
	"product_storage/rimport"
	"product_storage/tools/oradb"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	r := require.New(t)

	config, err := config.NewConfig(os.Getenv("CONF_PATH"))
	r.NoError(err)

	db := oradb.NewOracleDB(config.OracleConnectString())
	if err := db.Ping(); err != nil {
		r.NoError(err)
	}

	tx, err := db.Beginx()
	r.NoError(err)
	defer tx.Rollback()

	log := NewFileLogger("test")

	i := rimport.NewRepositoryImports(transaction.NewSQLSessionManager(db))

	dblog := NewDBLog("test", i)

	log.WithField("contractID", 3378).Infoln("InfoLn")
	log.WithField("oper_login", "m.zarif@sarkor.uz").Debugln("Debugln")
	log.WithFields(logrus.Fields{"contractID": 3378, "se_id": 200189621, "oper_login": "m.zarif@sarkor.uz", "details": "test"}).Errorln("Errorln")

	dblog.WithField("c_id", 3378).Infoln("InfoLn")
	dblog.Infoln("no field test")
	dblog.WithField("oper_login", "m.zarif@sarkor.uz").Debugln("Debugln")
	dblog.WithFields(logrus.Fields{"c_id": 3378, "se_id": 200189621, "oper_login": "m.zarif@sarkor.uz", "details": "test"}).Errorln("Errorln")
	dblog.WithFields(logrus.Fields{"c_id": 3378, "se_id": 200189621, "oper_login": "m.zarif@sarkor.uz", "details": "test"}).Infoln("All fields")

}

func TestUsecaseLogger(t *testing.T) {
	r := require.New(t)

	config, err := config.NewConfig(os.Getenv("CONF_PATH"))
	r.NoError(err)

	db := oradb.NewOracleDB(config.OracleConnectString())
	if err := db.Ping(); err != nil {
		r.NoError(err)
	}

	log := NewFileLogger("test_module")

	i := rimport.NewRepositoryImports(transaction.NewSQLSessionManager(db))

	ulog := NewUsecaseLogger(log, "myusecase")
	ulog.WithField("contractID", 3378).Infoln("InfoLn")
	ulog.WithField("oper_login", "m.zarif@sarkor.uz").Debugln("Debugln")
	ulog.WithFields(logrus.Fields{"contractID": 3378, "se_id": 200189621, "oper_login": "m.zarif@sarkor.uz", "details": "test"}).Errorln("Errorln")

	dblog := NewDBLog("test_module", i)
	udblog := NewUsecaseLogger(dblog, "myusecase")

	udblog.WithField("c_id", 3378).Infoln("InfoLn")

	udblog.Infoln("no field test")
	udblog.WithField("oper_login", "m.zarif@sarkor.uz").Debugln("Debugln")
	udblog.WithFields(logrus.Fields{"c_id": 3378, "se_id": 200189621, "oper_login": "m.zarif@sarkor.uz", "details": "test"}).Errorln("Errorln")
	udblog.WithFields(logrus.Fields{"c_id": 3378, "se_id": 200189621, "oper_login": "m.zarif@sarkor.uz", "details": "test"}).Infoln("All fields")

}

// func BenchmarkLogger(t *testing.B) {
// 	log := NewFileLogger("test_module")
// 	t.ResetTimer()

// 	for i := 0; i < t.N; i++ {
// 		log.Infoln("InfoLn")
// 	}
// }

// func BenchmarkLoggerWithField(t *testing.B) {
// 	log := NewFileLogger("test_module")

// 	t.ResetTimer()

// 	for i := 0; i < t.N; i++ {
// 		log.WithField("contractID", 3378).Infoln("InfoLn")
// 	}
// }

// func BenchmarkUsecaseLogger(t *testing.B) {
// 	log := NewNoFileLogger("test_module")
// 	ulog := NewUsecaseLogger(log, "myusecase")

// 	t.ResetTimer()

// 	for i := 0; i < t.N; i++ {
// 		ulog.Infoln("InfoLn")
// 	}
// }

// func BenchmarkUsecaseLoggerWithField(t *testing.B) {
// 	log := NewFileLogger("test_module")
// 	os.Stdout, _ = os.OpenFile(os.DevNull, os.O_WRONLY, 0)

// 	ulog := NewUsecaseLogger(log, "myusecase")

// 	t.ResetTimer()

// 	for i := 0; i < t.N; i++ {
// 		ulog.WithField("contractID", 3378).Infoln("InfoLn")
// 	}
// }
