package service

import (
	"errors"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/Rhymond/go-money"
)

const (
	dateLayout = "2006.01.02 15:04:05"
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

type storage interface {
	AddOrder(models.Order)
	RemoveOrder(int)
	UpdateOrder(int, models.Order)
	GetByID(int) models.Order
	GetByUserID(int) []models.Order
	GetReturns() []models.Order
	OrderHistory() []models.Order
	Save() error
	Contains(int) bool
}

type OrderService struct {
	Storage storage
}

func (s *OrderService) pack(order *models.Order, packaging models.Packaging) error {
	if packaging.GetCheckWeight() && order.Weight < packaging.GetMinWeight() {
		return errNotEnoughWeight
	}

	if order.Packaging == models.WrapPackaging || order.ExtraPackaging != models.WrapPackaging {
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

	if s.Storage.Contains(orderID) {
		return errOrderAlreadyExists
	}

	if weight < 0 {
		return errWrongWeight
	}

	if ok, _ := price.GreaterThan(money.New(0, money.RUB)); !ok {
		return errWrongPrice
	}

	currentTime := time.Now().Format(dateLayout)

	currentOrder := *models.NewOrder(orderID, userID, weight, price, orderStored,
		currentTime, expiryDate.Format(dateLayout), currentTime)

	for _, somePackaging := range packagings {
		err := s.pack(&currentOrder, somePackaging)
		if err != nil {
			return err
		}
	}

	s.Storage.AddOrder(currentOrder)

	return nil
}

func (s *OrderService) AcceptOrders(orders map[string]models.Order) int {
	ordersFailed := 0

	for _, someOrder := range orders {
		if s.Storage.Contains(someOrder.ID) {
			ordersFailed++

			continue
		}
		s.Storage.AddOrder(someOrder)
	}

	return ordersFailed
}

func (s *OrderService) ReturnOrder(orderID int) error {
	if !s.Storage.Contains(orderID) {
		return errOrderNotFound
	}
	someOrder := s.Storage.GetByID(orderID)
	if someOrder.Status == orderGiven {
		return errOrderIsGiven
	}
	date, _ := time.Parse(dateLayout, someOrder.ExpiryDate)
	if !date.Before(time.Now()) {
		return errOrderIsNotExpired
	}

	s.Storage.RemoveOrder(orderID)

	return nil
}

func isBeforeDeadline(someOrder models.Order, action string) bool {
	date := time.Now()
	var deadline time.Time
	switch action {
	case returnOrder:
		deadline, _ = time.Parse(dateLayout, someOrder.LastChange)

		return date.After(deadline)
	case giveOrder:
		deadline, _ = time.Parse(dateLayout, someOrder.ExpiryDate)

		return date.Before(deadline)
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
	if !s.Storage.Contains(orderID) {
		return errOrderNotFound
	}
	someOrder := s.Storage.GetByID(orderID)
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
	someOrder.LastChange = time.Now().Format(dateLayout)
	s.Storage.UpdateOrder(orderID, someOrder)

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
	orders := s.Storage.GetByUserID(userID)

	if count == 0 {
		return orders
	}

	return orders[:count]
}

func (s *OrderService) Returns() []models.Order {
	return s.Storage.GetReturns()
}

func (s *OrderService) OrderHistory() []models.Order {
	return s.Storage.OrderHistory()
}

func (s *OrderService) Save() error {
	return s.Storage.Save()
}
