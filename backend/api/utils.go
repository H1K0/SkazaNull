package api

import (
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
)

func HandleDBError(err error) (int, string) {
	pgErr, ok := err.(*pgconn.PgError)
	if !ok {
		return 500, err.Error()
	}
	statusStr := pgErr.Message[:3]
	message := pgErr.Message[4:]
	status, err := strconv.ParseInt(statusStr, 10, 0)
	if err == nil {
		return int(status), message
	} else {
		return 500, message
	}
}
