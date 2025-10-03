package apperror

import (
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
)

const duplicateErrCode = 1062

var mysqlDuplicateErr = &mysql.MySQLError{Number: duplicateErrCode}

func ToAppError(err error) error {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	case mysqlDuplicateErr.Is(err):
		return ErrDuplicateEntity
	default:
		return err
	}
}
