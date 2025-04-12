package order

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

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
		span, ctx := opentracing.StartSpanFromContext(ctx, "service.ProcessOrder")
		defer span.Finish()

		if ok, err := s.Storage.Contains(ctx, tx, orderID); err != nil || !ok {
			s.logger.Error(ErrOrderNotFound.Error(),
				zap.Int("id", orderID),
				zap.Int("user_id", userID),
				zap.String("action", action),
				zap.Error(ErrOrderNotFound),
			)
			span.SetTag("error", ErrOrderNotFound)

			return ErrOrderNotFound
		}

		someOrder, err := s.Storage.GetByID(ctx, tx, orderID)
		if err != nil {
			span.SetTag("error", err)

			return err
		}

		if !isOrderEligible(someOrder, userID, action) {
			s.logger.Error(ErrOrderNotEligible.Error(),
				zap.Int("id", orderID),
				zap.Int("user_id", userID),
				zap.String("action", action),
				zap.Error(ErrOrderNotEligible),
			)
			span.SetTag("error", ErrOrderNotEligible)

			return ErrOrderNotEligible
		}

		switch action {
		case giveOrder:
			someOrder.Status = models.GivenOrder
		case returnOrder:
			someOrder.Status = models.ReturnedOrder
		default:
			s.logger.Error(ErrUndefinedAction.Error(),
				zap.String("action", action),
				zap.Error(ErrUndefinedAction),
			)
			span.SetTag("error", ErrUndefinedAction)

			return ErrUndefinedAction
		}

		someOrder.LastChange = time.Now()

		return s.Storage.UpdateOrder(ctx, tx, orderID, someOrder)
	})
}
