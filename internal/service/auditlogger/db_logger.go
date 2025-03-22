package auditlogger

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) dbLogger(ctx context.Context, batches <-chan models.Log) {
	go func() {
		for batch := range batches {
			s.Storage.CreateLog(ctx, batch)
		}
	}()
}
