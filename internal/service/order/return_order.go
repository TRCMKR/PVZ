package order

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// ReturnOrder returns order by id
func (s *Service) ReturnOrder(ctx context.Context, orderID int) error {
	return s.txManager.RunSerializable(ctx, func(ctx context.Context, tx pgx.Tx) error {
		span, ctx := opentracing.StartSpanFromContext(ctx, "service.ReturnOrder")
		defer span.Finish()

		if ok, err := s.Storage.Contains(ctx, tx, orderID); err != nil || !ok {
			s.logger.Error(ErrOrderNotFound.Error(),
				zap.Int("id", orderID),
				zap.Error(err),
			)
			span.SetTag("error", ErrOrderNotFound)

			return ErrOrderNotFound
		}

		someOrder, err := s.Storage.GetByID(ctx, tx, orderID)
		if err != nil {
			span.SetTag("error", err)

			return err
		}

		if someOrder.Status == models.GivenOrder {
			s.logger.Error(ErrOrderIsGiven.Error(),
				zap.Int("id", orderID),
				zap.Error(ErrOrderIsGiven),
			)
			span.SetTag("error", ErrOrderIsGiven)

			return ErrOrderIsGiven
		}

		if !someOrder.ExpiryDate.Before(time.Now()) {
			s.logger.Error(ErrOrderIsNotExpired.Error(),
				zap.Int("id", orderID),
				zap.Error(ErrOrderIsNotExpired),
			)
			span.SetTag("error", ErrOrderIsNotExpired)

			return ErrOrderIsNotExpired
		}

		return s.Storage.RemoveOrder(ctx, tx, orderID)
	})
}
