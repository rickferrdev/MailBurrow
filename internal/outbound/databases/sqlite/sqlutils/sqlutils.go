package sqlutils

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"

	"github.com/rickferrdev/mail-burrow/internal/app/ports"
	"modernc.org/sqlite"
)

func NewError(err error) error {
	if err == nil {
		return nil
	}

	var portErr *ports.Error
	if errors.As(err, &portErr) {
		return portErr
	}

	if errors.Is(err, sql.ErrNoRows) {
		return ports.NewDatabaseRowsNotFound(err)
	}

	if errors.Is(err, context.Canceled) ||
		errors.Is(err, context.DeadlineExceeded) {
		return newDatabaseQuery(err)
	}

	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		return newSQLiteError(err, sqliteErr)
	}

	if mapped := fallbackSQLiteError(err); mapped != nil {
		return mapped
	}

	return newDatabaseQuery(err)
}

func NewResultError(result sql.Result) error {
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewError(err)
	}

	if rowsAffected == 0 {
		return ports.NewDatabaseRowsNotFound(nil)
	}

	return nil
}

func newSQLiteError(trace error, sqliteErr *sqlite.Error) error {
	switch sqliteErr.Code() {
	case sqliteConstraintUnique,
		sqliteConstraintPrimaryKey:
		return ports.NewDatabaseDuplicate(trace)

	case sqliteConstraintForeignKey:
		return ports.NewDatabaseForeignKey(trace)

	case sqliteConstraint,
		sqliteConstraintNotNull,
		sqliteConstraintCheck:
		return ports.NewDatabaseConstraint(trace)

	case sqliteBusy:
		return ports.NewDatabaseBusy(trace)

	case sqliteLocked:
		return ports.NewDatabaseLocked(trace)

	case sqliteReadonly:
		return ports.NewDatabaseReadonly(trace)

	case sqliteCantOpen:
		return ports.NewDatabaseConnection(trace)

	case sqliteRange,
		sqliteMismatch:
		return ports.NewDatabaseInvalidArgument(trace)

	case sqliteMisuse:
		return ports.NewDatabaseSQLiteDriver(trace)

	default:
		return newDatabaseQuery(trace)
	}
}

func fallbackSQLiteError(err error) error {
	msg := strings.ToLower(err.Error())

	switch {
	case strings.Contains(msg, "sql: no rows in result set"),
		strings.Contains(msg, "no rows in result set"):
		return ports.NewDatabaseRowsNotFound(err)

	case strings.Contains(msg, "unique constraint failed"),
		strings.Contains(msg, "primary key"),
		strings.Contains(msg, "constraint failed") && strings.Contains(msg, "unique"):
		return ports.NewDatabaseDuplicate(err)

	case strings.Contains(msg, "foreign key constraint failed"),
		strings.Contains(msg, "foreign key"):
		return ports.NewDatabaseForeignKey(err)

	case strings.Contains(msg, "not null constraint failed"),
		strings.Contains(msg, "check constraint failed"),
		strings.Contains(msg, "constraint failed"),
		strings.Contains(msg, "constraint violation"):
		return ports.NewDatabaseConstraint(err)

	case strings.Contains(msg, "database is busy"),
		strings.Contains(msg, "database busy"):
		return ports.NewDatabaseBusy(err)

	case strings.Contains(msg, "database is locked"),
		strings.Contains(msg, "database table is locked"),
		strings.Contains(msg, "database locked"):
		return ports.NewDatabaseLocked(err)

	case strings.Contains(msg, "attempt to write a readonly database"),
		strings.Contains(msg, "readonly database"),
		strings.Contains(msg, "database is readonly"):
		return ports.NewDatabaseReadonly(err)

	case strings.Contains(msg, "unable to open database file"):
		return ports.NewDatabaseConnection(err)

	case strings.Contains(msg, "invalid argument"),
		strings.Contains(msg, "unsupported type"),
		strings.Contains(msg, "converting argument"):
		return ports.NewDatabaseInvalidArgument(err)

	default:
		return nil
	}
}

func newDatabaseQuery(trace error) *ports.Error {
	return ports.New(
		trace,
		ports.MessageDatabaseQuery,
		ports.CodeDatabaseQuery,
		http.StatusInternalServerError,
	)
}

const (
	sqliteBusy       = 5
	sqliteLocked     = 6
	sqliteReadonly   = 8
	sqliteCantOpen   = 14
	sqliteConstraint = 19
	sqliteMismatch   = 20
	sqliteMisuse     = 21
	sqliteRange      = 25
)

const (
	sqliteConstraintCheck      = sqliteConstraint | 1<<8 // 275
	sqliteConstraintForeignKey = sqliteConstraint | 3<<8 // 787
	sqliteConstraintNotNull    = sqliteConstraint | 5<<8 // 1299
	sqliteConstraintPrimaryKey = sqliteConstraint | 6<<8 // 1555
	sqliteConstraintUnique     = sqliteConstraint | 8<<8 // 2067
)
