package order

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// UserOrders gets all order from this user
func (s *Service) UserOrders(ctx context.Context, userID int, count int) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.UserOrders")
	defer span.Finish()

	return s.Storage.GetByUserID(ctx, nil, userID, count)
}
