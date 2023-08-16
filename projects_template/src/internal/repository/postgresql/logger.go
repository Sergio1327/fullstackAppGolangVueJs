package postgresql

import (
	"product_storage/internal/entity/log"
	"product_storage/internal/repository"
	"product_storage/internal/transaction"
)

type loggerRepository struct{}

func NewLoggerRepository() repository.Logger {
	return &loggerRepository{}
}

func (l *loggerRepository) SaveLog(
	ts transaction.Session,
	row log.Row) error {
	sqlQuery := `
	insert into log_table 
	( logtime, flag, msg, module, fl, ln ) 
	values ( $1, $2, $3, $4, $5, $6)`

	_, err := SqlxTx(ts).Exec(sqlQuery,
		row.Time,
		row.Flag,
		row.Message,
		row.Module,
		row.File,
		row.Line,
	)
	return err
}

func (l *loggerRepository) SaveLogWithReturnID(
	ts transaction.Session,
	row log.Row,) (logID int, err error) {

	sqlQuery := `
	insert into log_table
	( logtime, flag, msg, module, fl, ln ) 
	values ( $1, $2, $3, $4, $5, $6 )
	returning log_id
	`
	err = SqlxTx(ts).QueryRow(sqlQuery,
		row.Time,
		row.Flag,
		row.Message,
		row.Module,
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
