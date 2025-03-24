package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

var (
	errCreateLog = errors.New("error creating log")
)

type LogsRepo struct {
	db database
}

func NewLogsRepo(db database) *LogsRepo {
	return &LogsRepo{
		db: db,
	}
}

func (r *LogsRepo) CreateLog(ctx context.Context, logBatch []models.Log) error {
	queryBatch := &pgx.Batch{}
	for _, log := range logBatch {
		queryBatch.Queue(`
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
	}

	br := r.db.SendBatch(ctx, queryBatch)
	defer br.Close()

	for i := 0; i < len(logBatch); i++ {
		if _, err := br.Exec(); err != nil {
			fmt.Println(3)
			fmt.Println(err.Error())
			return err
		}
	}

	return nil
}
