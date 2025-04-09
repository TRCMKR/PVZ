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
	errCreateJob = errors.New("error creating job")
)

// LogsRepo is a repository for logs table
type LogsRepo struct {
	db database
}

// NewLogsRepo creates instance of a LogsRepo
func NewLogsRepo(db database) *LogsRepo {
	return &LogsRepo{
		db: db,
	}
}

// CreateLog creates log
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
			log.OrderID, log.AdminID, log.Message, log.Date, log.URL, log.Method, log.Status)
	}

	br := r.db.SendBatch(ctx, queryBatch)
	defer br.Close()

	for i := 0; i < len(logBatch); i++ {
		if _, err := br.Exec(); err != nil {
			fmt.Println(err.Error())

			return errCreateLog
		}
	}

	return nil
}

// CreateJob creates job from logs in batches
func (r *LogsRepo) CreateJob(ctx context.Context, logBatch []models.Log) error {
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
												 status,
								                 job_status,
								                 attempts_left,
								                 updated_at)
								VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
								`,
			log.OrderID, log.AdminID, log.Message, log.Date, log.URL, log.Method, log.Status,
			log.JobStatus, log.AttemptsLeft, log.UpdatedAt)
	}

	br := r.db.SendBatch(ctx, queryBatch)
	defer br.Close()

	for i := 0; i < len(logBatch); i++ {
		if _, err := br.Exec(); err != nil {
			fmt.Println(err.Error())

			return errCreateJob
		}
	}

	return nil
}

// GetAndMarkLogs gets logs and marks them as being processed
func (r *LogsRepo) GetAndMarkLogs(ctx context.Context, batchSize int) ([]models.Log, error) {
	logs := make([]models.Log, 0)
	err := r.db.Select(ctx, &logs, `
									WITH cte AS (
										SELECT id
										FROM logs
										WHERE job_status = 1 OR job_status = 3 OR
											(job_status = 2 AND updated_at < (now() - INTERVAL '5 minutes'))
										ORDER BY date
										LIMIT $1
									),
									updated_logs AS (
										UPDATE logs
											SET job_status = 2,
												updated_at = now()
											WHERE id IN (SELECT id FROM cte)
											RETURNING *
									)
									SELECT * FROM updated_logs ORDER BY date;
									`, batchSize)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// GetLogs returns all logs in db
func (r *LogsRepo) GetLogs(ctx context.Context) ([]models.Log, error) {
	logs := make([]models.Log, 0)

	err := r.db.Select(ctx, &logs, `SELECT * FROM logs;`)
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// UpdateLog updates logs status and attempts left count
func (r *LogsRepo) UpdateLog(ctx context.Context, id int, newStatus int, attemptsLeft int) error {
	_, err := r.db.Exec(ctx, `UPDATE logs SET attempts_left = $1, job_status = $2 WHERE id = $3`,
		attemptsLeft, newStatus, id)

	return err
}
