package repository

import (
	"context"
	"errors"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
)

var (
	errCreateLog = errors.New("error creating log")
)

type LogsRepo struct {
	db database
}

func NewLogsRepo(db postgres.Database) *LogsRepo {
	return &LogsRepo{
		db: db,
	}
}

func (r *LogsRepo) CreateLog(ctx context.Context, log models.Log) error {
	_, err := r.db.Exec(ctx, `
							INSERT INTO logs(
							                 order_id,
							                 admin_id,
							                 message,
							                 date,
							                 url,
							                 method,
							                 status)
							VALUES ($1, $2, $3, $4, $5, $6, $7)
							`,
		log.OrderID, log.AdminID, log.Message, log.Date, log.Url, log.Method, log.Status)

	if err != nil {
		return errCreateLog
	}

	return nil
}
