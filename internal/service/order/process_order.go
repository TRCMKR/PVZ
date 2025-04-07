package order

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func isBeforeDeadline(someOrder models.Order, action string) bool {
	date := time.Now()
	switch action {
	case returnOrder:
		return date.After(someOrder.LastChange)
	case giveOrder:
		return date.Before(someOrder.ExpiryDate)
	}

	return false
}

func isOrderEligible(order models.Order, userID int, action string) bool {
	if !isBeforeDeadline(order, action) || order.UserID != userID {
		return false
	}

	if action == returnOrder {
		return order.Status == models.GivenOrder
	}

	return order.Status == models.StoredOrder
}

// ProcessOrder gives/returns order
func (s *Service) ProcessOrder(ctx context.Context, userID int, orderID int, action string) error {
	return s.txManager.RunSerializable(ctx, func(ctx context.Context, tx pgx.Tx) error {
		if ok, err := s.Storage.Contains(ctx, tx, orderID); err != nil || !ok {
			return ErrOrderNotFound
		}

		someOrder, err := s.Storage.GetByID(ctx, tx, orderID)
		if err != nil {
			return err
		}

		if !isOrderEligible(someOrder, userID, action) {
			return ErrOrderNotEligible
		}

		switch action {
		case giveOrder:
			someOrder.Status = models.GivenOrder
		case returnOrder:
			someOrder.Status = models.ReturnedOrder
		default:
			return ErrUndefinedAction
		}

		someOrder.LastChange = time.Now()

		return s.Storage.UpdateOrder(ctx, tx, orderID, someOrder)
	})
}
