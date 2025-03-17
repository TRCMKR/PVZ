package order

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) AcceptOrders(ctx context.Context, orders map[string]models.Order) (int, error) {
	ordersFailed := 0

	for _, someOrder := range orders {
		if ok, err := s.Storage.Contains(ctx, someOrder.ID); err != nil || !ok {
			ordersFailed++

			continue
		}
		err := s.Storage.AddOrder(ctx, someOrder)
		if err != nil {
			return ordersFailed, err
		}
	}

	return ordersFailed, nil
}
