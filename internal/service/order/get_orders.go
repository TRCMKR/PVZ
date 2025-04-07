package order

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
)

// GetOrders gets orders that satisfy conditions
func (s *Service) GetOrders(ctx context.Context, conds []query.Cond, count int, page int) ([]models.Order, error) {
	return s.Storage.GetOrders(ctx, nil, conds, count, page)
}
