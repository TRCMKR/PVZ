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

	if order.Packaging == models.WrapPackaging || order.ExtraPackaging != models.NoPackaging {
		return ErrWrongPackaging
	}

	if order.Packaging == models.NoPackaging {
		order.Packaging = packaging.GetType()
	} else {
		if packaging.GetType() != models.WrapPackaging {
			return ErrWrongPackaging
		}
		order.ExtraPackaging = packaging.GetType()
	}

	tmp, err := order.Price.Add(packaging.GetCost())
	if err != nil {
		return err
	}
	order.Price = *tmp

	return nil
}

func (s *Service) AcceptOrder(ctx context.Context, orderID int, userID int, weight float64, price money.Money,
	expiryDate time.Time, packagings []models.Packaging) error {
	if expiryDate.Before(time.Now()) {
		return ErrOrderExpired
	}

	var ok bool
	var err error
	if ok, err = s.Storage.Contains(ctx, orderID); ok {
		return ErrOrderAlreadyExists
	}
	if err != nil {
		return err
	}

	if weight < 0 {
		return ErrWrongWeight
	}

	if ok, err = price.GreaterThan(money.New(0, money.RUB)); err != nil || !ok {
		return ErrWrongPrice
	}

	currentTime := time.Now()

	currentOrder := *models.NewOrder(orderID, userID, weight, price, models.StoredOrder,
		currentTime, expiryDate, currentTime)

	for _, somePackaging := range packagings {
		if somePackaging == nil {
			return ErrUndefinedPackaging
		}
		err = s.pack(&currentOrder, somePackaging)
		if err != nil {
			return err
		}
	}

	err = s.Storage.AddOrder(ctx, currentOrder)
	if err != nil {
		return err
	}

	return nil
}
