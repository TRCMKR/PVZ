package order

import (
	"context"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/Rhymond/go-money"
)

func (s *Service) pack(order *models.Order, packaging models.Packaging) error {
	if packaging.GetCheckWeight() && order.Weight < packaging.GetMinWeight() {
		return ErrNotEnoughWeight
	}

	if order.Packaging == models.NoPackaging {
		order.Packaging = packaging.GetType()
	} else {
		order.ExtraPackaging = packaging.GetType()
	}

	tmp, err := order.Price.Add(packaging.GetCost())
	if err != nil {
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
		return ErrWrongPackaging
	}

	return nil
}

func (s *Service) validateOrder(order models.Order) error {
	if order.ExpiryDate.Before(time.Now()) {
		return ErrOrderExpired
	}

	if order.Weight < 0 {
		return ErrWrongWeight
	}

	if ok, err := order.Price.GreaterThan(money.New(0, money.RUB)); err != nil || !ok {
		return ErrWrongPrice
	}

	return nil
}

func (s *Service) AcceptOrder(ctx context.Context, orderID int, userID int, weight float64, price money.Money,
	expiryDate time.Time, packagings []models.Packaging) error {
	err := s.checkPackaging(packagings[0].GetType(), packagings[1].GetType())
	if err != nil {
		return err
	}

	currentTime := time.Now()

	currentOrder := *models.NewOrder(orderID, userID, weight, price, models.StoredOrder,
		currentTime, expiryDate, currentTime)

	err = s.validateOrder(currentOrder)
	if err != nil {
		return err
	}

	for _, somePackaging := range packagings {
		err = s.pack(&currentOrder, somePackaging)
		if err != nil {
			return err
		}
	}

	return s.txManager.RunReadCommitted(ctx, func(ctx context.Context) error {
		if ok, err := s.Storage.Contains(ctx, currentOrder.ID); ok {
			return ErrOrderAlreadyExists
		} else if err != nil {
			return err
		}

		return s.Storage.AddOrder(ctx, currentOrder)
	})
}
