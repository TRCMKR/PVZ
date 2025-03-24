package auditlogger

import (
	"context"
	"errors"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

var (
	errWritingLog = errors.New("failed to write audit log")
)

func (s *Service) dbLogger(ctx context.Context, batches <-chan []models.Log) error {
	for batch := range batches {
		select {
		case <-ctx.Done():
			timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			err := s.Storage.CreateLog(timeoutCtx, batch)
			if err != nil {
				return errWritingLog
			}

			return nil
		default:
			err := s.Storage.CreateLog(ctx, batch)
			if err != nil {
				return errWritingLog
			}
		}
	}

	return nil
}
