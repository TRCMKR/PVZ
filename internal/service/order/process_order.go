package order

import (
	"context"
	"time"

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

// ProcessOrder ...
func (s *Service) ProcessOrder(ctx context.Context, userID int, orderID int, action string) error {
	return s.txManager.RunSerializable(ctx, func(ctx context.Context) error {
		if ok, err := s.Storage.Contains(ctx, orderID); err != nil || !ok {
			return ErrOrderNotFound
		}

		someOrder, err := s.Storage.GetByID(ctx, orderID)
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

		return s.Storage.UpdateOrder(ctx, orderID, someOrder)
	})
}
