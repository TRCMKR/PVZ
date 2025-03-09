package service

import (
	"context"
	"errors"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/Rhymond/go-money"
)

const (
	giveOrder   = "give"
	returnOrder = "return"
)

const (
	orderGiven    = "given"
	orderReturned = "returned"
	orderStored   = "stored"
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
)

type orderStorage interface {
	AddOrder(context.Context, models.Order)
	RemoveOrder(context.Context, int)
	UpdateOrder(context.Context, int, models.Order)
	GetByID(context.Context, int) models.Order
	GetByUserID(context.Context, int, int) []models.Order
	GetReturns(context.Context) []models.Order
	GetOrders(context.Context, map[string]string, int, int) []models.Order
	Save(context.Context) error
	Contains(context.Context, int) bool
}

type OrderService struct {
	Storage orderStorage
	Ctx     context.Context
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

func (s *OrderService) AcceptOrder(orderID int, userID int, weight float64, price money.Money, expiryDate time.Time,
	packagings []models.Packaging) error {
	if expiryDate.Before(time.Now()) {
		return errOrderExpired
	}

	if s.Storage.Contains(s.Ctx, orderID) {
		return errOrderAlreadyExists
	}

	if weight < 0 {
		return errWrongWeight
	}

	if ok, _ := price.GreaterThan(money.New(0, money.RUB)); !ok {
		return errWrongPrice
	}

	currentTime := time.Now()

	currentOrder := *models.NewOrder(orderID, userID, weight, price, orderStored,
		currentTime, expiryDate, currentTime)

	for _, somePackaging := range packagings {
		err := s.pack(&currentOrder, somePackaging)
		if err != nil {
			return err
		}
	}

	s.Storage.AddOrder(s.Ctx, currentOrder)

	return nil
}

func (s *OrderService) AcceptOrders(orders map[string]models.Order) int {
	ordersFailed := 0

	for _, someOrder := range orders {
		if s.Storage.Contains(s.Ctx, someOrder.ID) {
			ordersFailed++

			continue
		}
		s.Storage.AddOrder(s.Ctx, someOrder)
	}

	return ordersFailed
}

func (s *OrderService) ReturnOrder(orderID int) error {
	if !s.Storage.Contains(s.Ctx, orderID) {
		return errOrderNotFound
	}
	someOrder := s.Storage.GetByID(s.Ctx, orderID)
	if someOrder.Status == orderGiven {
		return errOrderIsGiven
	}
	if !someOrder.ExpiryDate.Before(time.Now()) {
		return errOrderIsNotExpired
	}

	s.Storage.RemoveOrder(s.Ctx, orderID)

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
		return order.Status == orderGiven
	}

	return order.Status == orderStored
}

func (s *OrderService) processOrder(userID int, orderID int, action string) error {
	if !s.Storage.Contains(s.Ctx, orderID) {
		return errOrderNotFound
	}
	someOrder := s.Storage.GetByID(s.Ctx, orderID)
	if !isOrderEligible(someOrder, userID, action) {
		return errOrderNotEligible
	}

	switch action {
	case giveOrder:
		someOrder.Status = orderGiven
	case returnOrder:
		someOrder.Status = orderReturned
	default:
		return errUndefinedAction
	}
	someOrder.LastChange = time.Now()
	s.Storage.UpdateOrder(s.Ctx, orderID, someOrder)

	return nil
}

func (s *OrderService) ProcessOrders(userID int, orderIDs []int, action string) (int, error) {
	ordersFailed := 0

	for _, orderID := range orderIDs {
		err := s.processOrder(userID, orderID, action)
		if err != nil {
			if errors.Is(err, errUndefinedAction) {
				return 0, errUndefinedAction
			}
			ordersFailed++
		}
	}

	return ordersFailed, nil
}

func (s *OrderService) UserOrders(userID int, count int) []models.Order {
	orders := s.Storage.GetByUserID(s.Ctx, userID, count)

	return orders
}

func (s *OrderService) Returns() []models.Order {
	return s.Storage.GetReturns(s.Ctx)
}

func (s *OrderService) Save() error {
	return s.Storage.Save(s.Ctx)
}

func (s *OrderService) GetOrders(params map[string]string, count int, page int) []models.Order {
	return s.Storage.GetOrders(s.Ctx, params, count, page)
}
