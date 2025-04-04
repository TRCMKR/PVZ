package order

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// ReturnOrder ...
func (s *Service) ReturnOrder(ctx context.Context, orderID int) error {
	return s.txManager.RunSerializable(ctx, func(ctx context.Context, tx pgx.Tx) error {
		if ok, err := s.Storage.Contains(ctx, tx, orderID); err != nil || !ok {
			return ErrOrderNotFound
		}

		someOrder, err := s.Storage.GetByID(ctx, tx, orderID)
		if err != nil {
			return err
		}

		if someOrder.Status == models.GivenOrder {
			return ErrOrderIsGiven
		}

		if !someOrder.ExpiryDate.Before(time.Now()) {
			return ErrOrderIsNotExpired
		}

		return s.Storage.RemoveOrder(ctx, tx, orderID)
	})
}
