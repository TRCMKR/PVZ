package auditlogger

import (
	"context"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) dbWorker(ctx context.Context, batchSize int, timeout time.Duration, jobs chan models.Log) error {
	batches := s.batcher(jobs, batchSize, timeout)
	return s.dbLogger(ctx, batches)
}

func (s *Service) stdoutWorker(ctx context.Context, batchSize int, timeout time.Duration,
	jobs chan models.Log, operator func(log models.Log) bool) error {
	filtered := s.filter(jobs, operator)
	batches := s.batcher(filtered, batchSize, timeout)
	return s.stdoutLogger(ctx, batches)
}
