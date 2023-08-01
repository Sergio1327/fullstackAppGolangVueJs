package oracle

import (
	"product_storage/internal/entity/log"
	"product_storage/internal/repository"
	"product_storage/internal/transaction"
	"product_storage/tools/sqlnull"

	"github.com/mattn/go-oci8"
)

type loggerRepository struct {
}

func NewLogger() repository.Logger {
	return &loggerRepository{}
}

func (r *loggerRepository) SaveLog(ts transaction.Session, row log.Row, contractID, seID sqlnull.NullInt64, operLogin sqlnull.NullString) error {
	sqlQuery := `insert into log_table (time, flag, msg, c_id, module, se_id, oper_login, caller, line_no) values (:1, :2, :3, :4, :5, :6, :7, :8, :9)`

	_, err := SqlxTx(ts).Exec(sqlQuery,
		row.Time,
		row.Flag,
		row.Message,
		contractID,
		row.Module,
		seID,
		operLogin,
		row.File,
		row.Line,
	)

	return err
}

func (r *loggerRepository) SaveLogWithReturnID(ts transaction.Session, row log.Row, contractID, seID sqlnull.NullInt64, operLogin sqlnull.NullString) (logID int, err error) {
	sqlQuery := `insert into log_table (time, flag, msg, c_id, module, se_id, oper_login, caller, line_no) values (:1, :2, :3, :4, :5, :6, :7, :8, :9)`

	res, err := SqlxTx(ts).Exec(sqlQuery,
		row.Time,
		row.Flag,
		row.Message,
		contractID,
		row.Module,
		seID,
		operLogin,
		row.File,
		row.Line,
	)
	if err != nil {
		return
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return
	}

	rowID := oci8.GetLastInsertId(lastInsertId)
	err = SqlxTx(ts).QueryRow("select log_id from log_table where rowid = :1", rowID).Scan(&logID)
	if err != nil {
		return
	}
	return
}

func (r *loggerRepository) SaveLogDetails(ts transaction.Session, logID int, details map[string]string) error {
	sqlQuery := `INSERT INTO log_details (LOG_ID, NAME, VALUE) VALUES (:1, :2, :3)`

	stmt, err := SqlxTx(ts).Prepare(sqlQuery)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for k, v := range details {
		if _, err := stmt.Exec(logID, k, v); err != nil {
			return err
		}
	}

	return nil
}
