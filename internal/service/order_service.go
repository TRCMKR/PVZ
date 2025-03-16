package service

import (
	"context"
	"errors"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"

	"github.com/Rhymond/go-money"
)

const (
	giveOrder   = "give"
	returnOrder = "return"
)

var (
	errOrderAlreadyExists = errors.New("such order exists")
	errOrderExpired       = errors.New("expired order")
	errOrderIsNotExpired  = errors.New("order is not expired")
	errOrderNotFound      = errors.New("order not found")
	errOrderIsGiven       = errors.New("order is given")
	errWrongWeight        = errors.New("wrong weight")
	errWrongPrice         = errors.New("wrong price")
	errOrderNotEligible   = errors.New("order not eligible")
	errUndefinedAction    = errors.New("undefined action")
	errNotEnoughWeight    = errors.New("not enough weight")
	errWrongPackaging     = errors.New("wrong packaging")
	errUndefinedPackaging = errors.New("undefined packaging")
)

type orderStorage interface {
	AddOrder(context.Context, models.Order) error
	RemoveOrder(context.Context, int) error
	UpdateOrder(context.Context, int, models.Order) error
	GetByID(context.Context, int) (models.Order, error)
	GetByUserID(context.Context, int, int) ([]models.Order, error)
	GetReturns(context.Context) ([]models.Order, error)
	GetOrders(context.Context, []query.Cond, int, int) ([]models.Order, error)
	Contains(context.Context, int) (bool, error)
}

type OrderService struct {
	Storage orderStorage
}

func (s *OrderService) pack(order *models.Order, packaging models.Packaging) error {
	if packaging.GetCheckWeight() && order.Weight < packaging.GetMinWeight() {
		return errNotEnoughWeight
	}

	if order.Packaging == models.WrapPackaging || order.ExtraPackaging != models.NoPackaging {
		return errWrongPackaging
	}

	if order.Packaging == models.NoPackaging {
		order.Packaging = packaging.GetType()
	} else {
		if packaging.GetType() != models.WrapPackaging {
			return errWrongPackaging
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

func (s *OrderService) AcceptOrder(ctx context.Context, orderID int, userID int, weight float64, price money.Money,
	expiryDate time.Time, packagings []models.Packaging) error {
	if expiryDate.Before(time.Now()) {
		return errOrderExpired
	}

	var ok bool
	var err error
	if ok, err = s.Storage.Contains(ctx, orderID); ok {
		return errOrderAlreadyExists
	}
	if err != nil {
		return err
	}

	if weight < 0 {
		return errWrongWeight
	}

	if ok, err = price.GreaterThan(money.New(0, money.RUB)); err != nil || !ok {
		return errWrongPrice
	}

	currentTime := time.Now()

	currentOrder := *models.NewOrder(orderID, userID, weight, price, models.StoredOrder,
		currentTime, expiryDate, currentTime)

	for _, somePackaging := range packagings {
		if somePackaging == nil {
			return errUndefinedPackaging
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

func (s *OrderService) AcceptOrders(ctx context.Context, orders map[string]models.Order) (int, error) {
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

func (s *OrderService) ReturnOrder(ctx context.Context, orderID int) error {
	if ok, err := s.Storage.Contains(ctx, orderID); err != nil || !ok {
		return errOrderNotFound
	}
	someOrder, err := s.Storage.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	if someOrder.Status == models.GivenOrder {
		return errOrderIsGiven
	}
	if !someOrder.ExpiryDate.Before(time.Now()) {
		return errOrderIsNotExpired
	}

	err = s.Storage.RemoveOrder(ctx, orderID)
	if err != nil {
		return err
	}

	return nil
}

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

func (s *OrderService) processOrder(ctx context.Context, userID int, orderID int, action string) error {
	if ok, err := s.Storage.Contains(ctx, orderID); err != nil || !ok {
		return errOrderNotFound
	}
	someOrder, err := s.Storage.GetByID(ctx, orderID)
	if err != nil {
		return err
	}
	if !isOrderEligible(someOrder, userID, action) {
		return errOrderNotEligible
	}

	switch action {
	case giveOrder:
		someOrder.Status = models.GivenOrder
	case returnOrder:
		someOrder.Status = models.ReturnedOrder
	default:
		return errUndefinedAction
	}
	someOrder.LastChange = time.Now()
	err = s.Storage.UpdateOrder(ctx, orderID, someOrder)
	if err != nil {
		return err
	}

	return nil
}

func (s *OrderService) ProcessOrders(ctx context.Context, userID int, orderIDs []int, action string) (int, error) {
	ordersFailed := 0

	for _, orderID := range orderIDs {
		err := s.processOrder(ctx, userID, orderID, action)
		if err != nil {
			if errors.Is(err, errUndefinedAction) {
				return 0, errUndefinedAction
			}
			ordersFailed++
		}
	}

	return ordersFailed, nil
}

func (s *OrderService) UserOrders(ctx context.Context, userID int, count int) ([]models.Order, error) {
	return s.Storage.GetByUserID(ctx, userID, count)
}

func (s *OrderService) Returns(ctx context.Context) ([]models.Order, error) {
	return s.Storage.GetReturns(ctx)
}

func (s *OrderService) GetOrders(ctx context.Context, conds []query.Cond, count int,
	page int) ([]models.Order, error) {
	return s.Storage.GetOrders(ctx, conds, count, page)
}
