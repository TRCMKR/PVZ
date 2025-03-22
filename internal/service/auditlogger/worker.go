package auditlogger

import (
	"context"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) dbWorker(ctx context.Context, batchSize int, timeout time.Duration, jobs chan models.Log) {
	batches := s.batcher(jobs, batchSize, timeout)
	s.dbLogger(ctx, batches)
}

func (s *Service) stduoutWorker(batchSize int, timeout time.Duration, jobs chan models.Log, operator func(log models.Log) bool) {
	filtered := s.filter(jobs, operator)
	batches := s.batcher(filtered, batchSize, timeout)
	s.stdoutLogger(batches)
}
