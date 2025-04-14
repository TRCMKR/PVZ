package order

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
)

// GetOrders gets orders that satisfy conditions
func (s *Service) GetOrders(ctx context.Context, conds []query.Cond, count int, page int) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.GetOrders")
	defer span.Finish()

	return s.Storage.GetOrders(ctx, nil, conds, count, page)
}
