package order

import (
	"context"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/Rhymond/go-money"
)

func (s *Service) pack(order *models.Order, packaging models.Packaging) error {
	if packaging.GetCheckWeight() && order.Weight < packaging.GetMinWeight() {
		s.logger.Error(ErrNotEnoughWeight.Error(),
			zap.Float64("weight", order.Weight),
			zap.Float64("min_weight", packaging.GetMinWeight()),
			zap.Int("packaging", int(packaging.GetType())),
			zap.Error(ErrNotEnoughWeight),
		)

		return ErrNotEnoughWeight
	}

	if order.Packaging == models.NoPackaging {
		order.Packaging = packaging.GetType()
	} else {
		order.ExtraPackaging = packaging.GetType()
	}

	tmp, err := order.Price.Add(packaging.GetCost())
	if err != nil {
		s.logger.Error("error adding price",
			zap.Int("packaging", int(packaging.GetType())),
			zap.Error(err),
		)

		return err
	}
	order.Price = *tmp

	return nil
}

func (s *Service) checkPackaging(packagingType models.PackagingType, extraPackagingType models.PackagingType) error {
	if ((packagingType == models.NoPackaging || packagingType == models.WrapPackaging) &&
		extraPackagingType != models.NoPackaging) ||
		(packagingType != models.NoPackaging && packagingType != models.WrapPackaging &&
			extraPackagingType != models.NoPackaging && extraPackagingType != models.WrapPackaging) {
		s.logger.Error(ErrWrongPackaging.Error(),
			zap.Int("packaging", int(packagingType)),
			zap.Int("extra_packaging", int(extraPackagingType)),
			zap.Error(ErrWrongPackaging),
		)

		return ErrWrongPackaging
	}

	return nil
}

func (s *Service) validateOrder(order models.Order) error {
	if order.ExpiryDate.Before(time.Now()) {
		s.logger.Error(ErrOrderExpired.Error(),
			zap.Int("order_id", order.ID),
			zap.Time("expiry_date", order.ExpiryDate),
			zap.Error(ErrOrderExpired),
		)

		return ErrOrderExpired
	}

	if order.Weight < 0 {
		s.logger.Error("negative weight",
			zap.Int("order_id", order.ID),
			zap.Float64("weight", order.Weight),
			zap.Error(ErrWrongWeight),
		)

		return ErrWrongWeight
	}

	if ok, err := order.Price.GreaterThan(money.New(0, money.RUB)); err != nil || !ok {
		s.logger.Error("negative price",
			zap.Int("order_id", order.ID),
			zap.Int64("price", order.Price.Amount()),
			zap.Error(ErrWrongPrice),
		)

		return ErrWrongPrice
	}

	return nil
}

// AcceptOrder accept order
func (s *Service) AcceptOrder(ctx context.Context, orderID int, userID int, weight float64, price money.Money,
	expiryDate time.Time, packagings []models.Packaging) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "service.AcceptOrder")
	defer span.Finish()

	err := s.checkPackaging(packagings[0].GetType(), packagings[1].GetType())
	if err != nil {
		span.SetTag("error", err)

		return err
	}

	currentTime := time.Now()

	currentOrder := *models.NewOrder(orderID, userID, weight, price, models.StoredOrder,
		currentTime, expiryDate, currentTime)

	err = s.validateOrder(currentOrder)
	if err != nil {
		span.SetTag("error", err)

		return err
	}

	for _, somePackaging := range packagings {
		err = s.pack(&currentOrder, somePackaging)
		if err != nil {
			span.SetTag("error", err)

			return err
		}
	}

	return s.txManager.RunRepeatableRead(ctx, func(ctx context.Context, tx pgx.Tx) error {
		if ok, err := s.Storage.Contains(ctx, tx, currentOrder.ID); ok {
			s.logger.Error(ErrOrderAlreadyExists.Error(),
				zap.Int("order_id", orderID),
				zap.Error(ErrOrderAlreadyExists),
			)
			span.SetTag("error", ErrOrderAlreadyExists)

			return ErrOrderAlreadyExists
		} else if err != nil {
			span.SetTag("error", err)

			return err
		}

		return s.Storage.AddOrder(ctx, tx, currentOrder)
	})
}
