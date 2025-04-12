package order

import (
	"context"

	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// AcceptOrders accepts orders
func (s *Service) AcceptOrders(ctx context.Context, orders map[string]models.Order) (int, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.AcceptOrders")
	defer span.Finish()

	ordersFailed := 0

	for _, someOrder := range orders {
		if ok, err := s.Storage.Contains(ctx, nil, someOrder.ID); err != nil || !ok {
			ordersFailed++

			continue
		}

		err := s.Storage.AddOrder(ctx, nil, someOrder)
		if err != nil {
			span.SetTag("error", err)

			return ordersFailed, err
		}
	}

	return ordersFailed, nil
}
