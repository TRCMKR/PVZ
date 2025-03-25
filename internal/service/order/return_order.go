package order

import (
	"context"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) ReturnOrder(ctx context.Context, orderID int) error {
	if ok, err := s.Storage.Contains(ctx, orderID); err != nil || !ok {
		return ErrOrderNotFound
	}

	someOrder, err := s.Storage.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if someOrder.Status == models.GivenOrder {
		return ErrOrderIsGiven
	}

	if !someOrder.ExpiryDate.Before(time.Now()) {
		return ErrOrderIsNotExpired
	}

	err = s.Storage.RemoveOrder(ctx, orderID)
	if err != nil {
		return err
	}

	return nil
}
