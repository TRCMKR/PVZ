package jsondata

import (
	"errors"
	"homework/packaging"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/bytedance/sonic"

	"homework/order"
)

var (
	errOrderAlreadyExists  = errors.New("such order exists")
	errOrderExpired        = errors.New("expired order")
	errOrderIsNotExpired   = errors.New("order is not expired")
	errWrongDateFormat     = errors.New("wrong date format")
	errDataNotMarshalled   = errors.New("data wasn't marshalled")
	errDataNotUnmarshalled = errors.New("data wasn't unmarshalled")
	errDataNotWritten      = errors.New("data wasn't written")
	errOrderNotFound       = errors.New("order not found")
	errWrongArgument       = errors.New("wrong argument")
	errFileNotOpened       = errors.New("file not opened")
	errOrderIsGiven        = errors.New("order is given")
	errWrongWeight         = errors.New("wrong weight")
	errWrongPrice          = errors.New("wrong price")
	// errUserNotExists       = errors.New("user not exists")
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

const (
	dateLayout             = "2006.01.02 15:04:05"
	inputDateAndTimeLayout = "2006.01.02-15:04:05"
	inputDateLayout        = "2006.01.02"
)

type Storage struct {
	data map[string]order.Order
	path string
}

func New(path string) (*Storage, error) {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("read json fail", err)
	}

	data := make(map[string]order.Order)

	err = sonic.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal("unmarshal json fail", err)
	}

	return &Storage{
		data: data,
		path: path,
	}, nil
}

func (s *Storage) Save() error {
	jsonData, err := sonic.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return errDataNotMarshalled
	}
	err = os.WriteFile(s.path, jsonData, 0600)
	if err != nil {
		return errDataNotWritten
	}

	return nil
}

func parseInputDate(expiryDate string) (time.Time, error) {
	date, err := time.Parse(inputDateAndTimeLayout, expiryDate)
	if err != nil {
		date, err = time.Parse(inputDateLayout, expiryDate)
		if err != nil {
			return time.Time{}, errWrongDateFormat
		}
	}

	return date, nil
}

func (s *Storage) AcceptOrder(orderID int, userID int, weight float64, price float64, expiryDate string,
	packagings []packaging.Packaging) error {
	date, err := parseInputDate(expiryDate)
	if err != nil {
		return err
	}

	if date.Before(time.Now()) {
		return errOrderExpired
	}

	strOrderID := strconv.Itoa(orderID)
	if _, ok := s.data[strOrderID]; ok {
		return errOrderAlreadyExists
	}

	if weight < 0 {
		return errWrongWeight
	}

	if price < 0 {
		return errWrongPrice
	}

	err = packaging.CheckPackaging(packagings)
	if err != nil {
		return err
	}

	currentTime := time.Now().Format(dateLayout)

	currentOrder := order.Order{
		ID:          orderID,
		UserID:      userID,
		Weight:      weight,
		Price:       price,
		ArrivalDate: currentTime,
		Status:      orderStored,
		ExpiryDate:  date.Format(dateLayout),
		LastChange:  currentTime,
	}

	for _, somePackaging := range packagings {
		err = somePackaging.Pack(&currentOrder)
		if err != nil {
			return err
		}
	}
	currentOrder.Packaging = packaging.FormPackagingString(packagings)

	s.data[strOrderID] = currentOrder

	return nil
}

func (s *Storage) AcceptOrders(path string) (int, int, error) {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, errFileNotOpened
	}

	data := make(map[string]order.Order)

	err = sonic.Unmarshal(jsonData, &data)
	if err != nil {
		return 0, 0, errDataNotUnmarshalled
	}

	ordersFailed := 0
	orderCount := 0
	for key, value := range data {
		orderCount++
		if _, ok := s.data[key]; !ok {
			s.data[key] = value

			continue
		}

		ordersFailed++
	}

	return ordersFailed, orderCount, nil
}

func (s *Storage) ReturnOrder(orderID int) error {
	strOrderID := strconv.Itoa(orderID)
	if _, ok := s.data[strOrderID]; !ok {
		return errOrderNotFound
	}
	if s.data[strOrderID].Status == orderGiven {
		return errOrderIsGiven
	}
	date, _ := time.Parse(dateLayout, s.data[strOrderID].ExpiryDate)
	if !date.Before(time.Now()) {
		return errOrderIsNotExpired
	}

	delete(s.data, strOrderID)

	return nil
}

func isBeforeDeadline(someOrder order.Order, action string) bool {
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

func isOrderEligible(order order.Order, userID int, action string) bool {
	if !isBeforeDeadline(order, action) || order.UserID != userID {
		return false
	}
	if action == returnOrder {
		return order.Status == orderGiven
	}

	return order.Status == orderStored
}

func (s *Storage) processOrder(orderID, userID int, action string) bool {
	strOrderID := strconv.Itoa(orderID)
	someOrder, exists := s.data[strOrderID]
	if !exists || !isOrderEligible(someOrder, userID, action) {
		return false
	}

	statusMap := map[string]string{
		giveOrder:   orderGiven,
		returnOrder: orderReturned,
	}

	someOrder.Status = statusMap[action]
	someOrder.LastChange = time.Now().Format(dateLayout)
	s.data[strOrderID] = someOrder

	return true
}

func (s *Storage) ProcessOrders(userID int, orderIDs []int, action string) (int, error) {
	if action != giveOrder && action != returnOrder {
		return 0, errWrongArgument
	}

	ordersFailed := 0
	for _, orderID := range orderIDs {
		if !s.processOrder(orderID, userID, action) {
			ordersFailed++
		}
	}

	return ordersFailed, nil
}

func (s *Storage) UserOrders(userID int, count int) ([]order.Order, error) {
	userOrders := make([]order.Order, 0)

	currentCount := 0
	orderHistory := s.OrderHistory()
	for i := range orderHistory {
		if userID != orderHistory[i].UserID {
			continue
		}
		userOrders = append(userOrders, orderHistory[i])
		currentCount++
		if currentCount == count {
			break
		}
	}

	return userOrders, nil
}

func (s *Storage) Returns() []order.Order {
	returns := make([]order.Order, 0)

	for i := range s.data {
		if s.data[i].Status != orderReturned {
			continue
		}
		returns = append(returns, s.data[i])
	}

	return returns
}

func (s *Storage) OrderHistory() []order.Order {
	orderHistory := make([]order.Order, 0)
	for i := range s.data {
		orderHistory = append(orderHistory, s.data[i])
	}
	sort.Slice(orderHistory, func(i, j int) bool {
		t1, _ := time.Parse(dateLayout, orderHistory[i].LastChange)
		t2, _ := time.Parse(dateLayout, orderHistory[j].LastChange)

		return t2.Before(t1)
	})

	return orderHistory
}
