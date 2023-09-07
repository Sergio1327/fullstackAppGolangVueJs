package test

// import (
// 	"product_storage/internal/entity/log"
// 	"product_storage/internal/usecase"
// 	"product_storage/rimport"
// 	"product_storage/tools/logger"
// 	"product_storage/tools/sqlnull"
// 	"testing"

// 	"github.com/golang/mock/gomock"
// 	"github.com/stretchr/testify/require"
// )

// func TestSaveLog(t *testing.T) {
// 	r := require.New(t)
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	i := rimport.NewTestRepositoryImports(ctrl)
// 	ts := i.MockSessionWithCommit()
// 	i.SessionManager.EXPECT().CreateSession().Return(ts).AnyTimes()

// 	t.Run("без деталей", func(t *testing.T) {
// 		row := log.NewTestRow(nil)

// 		var (
// 			contractID sqlnull.NullInt64
// 			seID       sqlnull.NullInt64
// 			operLogin  sqlnull.NullString
// 		)

// 		contractID.Scan(row.SpecialFields["c_id"].Value)
// 		seID.Scan(row.SpecialFields["se_id"].Value)
// 		operLogin.Scan(row.SpecialFields["oper_login"].Value)

// 		i.MockRepository.Logger.EXPECT().SaveLog(
// 			ts,
// 			row,
// 			contractID,
// 			seID,
// 			operLogin,
// 		).Return(nil)

// 		u := usecase.NewLogger(logger.NewNoFileLogger("test"), i.RepositoryImports())

// 		err := u.SaveLog(row)
// 		r.NoError(err)

// 	})

// 	t.Run("с деталями", func(t *testing.T) {
// 		u := usecase.NewLogger(logger.NewNoFileLogger("test"), i.RepositoryImports())

// 		row := log.NewTestRow(map[string]string{
// 			"test1": "test2",
// 		})

// 		var (
// 			contractID sqlnull.NullInt64
// 			seID       sqlnull.NullInt64
// 			operLogin  sqlnull.NullString
// 		)

// 		contractID.Scan(row.SpecialFields["c_id"].Value)
// 		seID.Scan(row.SpecialFields["se_id"].Value)
// 		operLogin.Scan(row.SpecialFields["oper_login"].Value)

// 		logID := 1

// 		i.MockRepository.Logger.EXPECT().SaveLogWithReturnID(
// 			ts,
// 			row,
// 			contractID,
// 			seID,
// 			operLogin,
// 		).Return(logID, nil)

// 		i.MockRepository.Logger.EXPECT().SaveLogDetails(ts, logID, row.Details).Return(nil)

// 		err := u.SaveLog(row)
// 		r.NoError(err)
// 	})

// }
