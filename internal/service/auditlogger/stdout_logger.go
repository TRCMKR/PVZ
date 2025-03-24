package auditlogger

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) stdoutLogger(ctx context.Context, batches <-chan []models.Log) error {
	for batch := range batches {
		select {
		case <-ctx.Done():
			return nil
		default:
			for _, log := range batch {
				fmt.Print(log.String())
			}
		}
	}

	return nil
}
