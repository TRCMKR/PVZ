package order

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) Returns(ctx context.Context) ([]models.Order, error) {
	return s.Storage.GetReturns(ctx)
}
