package order

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// UserOrders ...
func (s *Service) UserOrders(ctx context.Context, userID int, count int) ([]models.Order, error) {
	return s.Storage.GetByUserID(ctx, userID, count)
}
