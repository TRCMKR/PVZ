package auditlogger

import (
	"context"
	"fmt"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) CreateLog(ctx context.Context, order models.Log) {
	select {
	case <-ctx.Done():
		fmt.Print("context canceled, writing log canceled")
		return
	default:
		s.jobs <- order
	}

}
