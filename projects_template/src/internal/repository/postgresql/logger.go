package postgresql

import (
	"product_storage/internal/entity/log"
	"product_storage/internal/repository"
	"product_storage/internal/transaction"
	"product_storage/tools/sqlnull"
)

type loggerRepository struct{}

func NewLoggerRepository() repository.Logger {
	return &loggerRepository{}
}

func (l *loggerRepository) SaveLog(
	ts transaction.Session,
	row log.Row,
	operLogin sqlnull.NullString) error {
	sqlQuery := `
	insert into log_table 
	( time, flag, msg, c_id, module, se_id, oper_login, caller, line_no ) 
	values ( $1, $2, $3, $4, $5, $6, $7 )`

	_, err := SqlxTx(ts).Exec(sqlQuery,
		row.Time,
		row.Flag,
		row.Message,
		row.Module,
		operLogin,
		row.File,
		row.Line,
	)
	return err
}

func (l *loggerRepository) SaveLogWithReturnID(
	ts transaction.Session,
	row log.Row,
	contractID,
	seID sqlnull.NullInt64,
	operLogin sqlnull.NullString) (logID int, err error) {
	sqlQuery := `
	insert into log_table
	( time, flag, msg, module, oper_login, caller, line_no ) values
	( $1, $2, $3, $4, $5, $6, $7 )
	returning log_id
	`
	err = SqlxTx(ts).QueryRow(sqlQuery,
		row.Time,
		row.Flag,
		row.Message,
		contractID,
		row.Module,
		seID,
		operLogin,
		row.File,
		row.Line,
	).Scan(&logID)

	return logID, err
}

func (l *loggerRepository) SaveLogDetails(ts transaction.Session, logID int, details map[string]string) error {
	sqlQuery := `
	INSERT INTO log_details 
	(LOG_ID, NAME, VALUE) 
	VALUES ($1, $2, $3)`

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
