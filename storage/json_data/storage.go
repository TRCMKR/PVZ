package json_data

import (
	"errors"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/bytedance/sonic"

	"homework/order"
)

var (
	ErrOrderAlreadyExists  = errors.New("such order exists")
	ErrOrderExpired        = errors.New("expired order")
	ErrOrderIsNotExpired   = errors.New("order is not expired")
	ErrWrongDateFormat     = errors.New("wrong date format")
	ErrDataNotMarshalled   = errors.New("data wasn't marshalled")
	ErrDataNotUnmarshalled = errors.New("data wasn't unmarshalled")
	ErrDataNotWritten      = errors.New("data wasn't written")
	ErrOrderNotFound       = errors.New("order not found")
	ErrWrongArgument       = errors.New("wrong argument")
	ErrFileNotOpened       = errors.New("file not opened")
	ErrOrderIsGiven        = errors.New("order is given")
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
		return ErrDataNotMarshalled
	}
	err = os.WriteFile(s.path, jsonData, 0600)
	if err != nil {
		return ErrDataNotWritten
	}

	return nil
}

func (s *Storage) AcceptOrder(orderID int, userID int, expiryDate string) error {
	date, err := time.Parse(inputDateAndTimeLayout, expiryDate)
	if err != nil {
		date, err = time.Parse(inputDateLayout, expiryDate)
		if err != nil {
			return ErrWrongDateFormat
		}
	}
	if date.Before(time.Now()) {
		return ErrOrderExpired
	}

	for i := range s.data {
		if _, ok := s.data[i]; ok {
			return ErrOrderAlreadyExists
		}
	}

	currentTime := time.Now().Format(dateLayout)

	s.data[strconv.Itoa(orderID)] = order.Order{
		ID:          orderID,
		UserID:      userID,
		ArrivalDate: currentTime,
		Status:      "stored",
		ExpiryDate:  date.Format(dateLayout),
		LastChange:  currentTime,
	}

	return nil
}

func (s *Storage) AcceptOrders(path string) (int, int, error) {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return 0, 0, ErrFileNotOpened
	}

	data := make(map[string]order.Order)

	err = sonic.Unmarshal(jsonData, &data)
	if err != nil {
		return 0, 0, ErrDataNotUnmarshalled
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
		return ErrOrderNotFound
	}
	if s.data[strOrderID].Status == "given" {
		return ErrOrderIsGiven
	}
	date, _ := time.Parse(dateLayout, s.data[strOrderID].ExpiryDate)
	if !date.Before(time.Now()) {
		return ErrOrderIsNotExpired
	}

	delete(s.data, strOrderID)

	return nil
}

func isBeforeDeadline(someOrder order.Order, action string) bool {
	date := time.Now()
	var deadline time.Time
	switch action {
	case "return":
		deadline, _ = time.Parse(dateLayout, someOrder.LastChange)
		return date.After(deadline)
	case "give":
		deadline, _ = time.Parse(dateLayout, someOrder.ExpiryDate)
		return date.Before(deadline)
	}

	return false
}

func isOrderEligible(order order.Order, action string) bool {
	if !isBeforeDeadline(order, action) {
		return false
	}
	if action == "return" {
		return order.Status == "given"
	}

	return order.Status == "stored"
}

func (s *Storage) ProcessOrders(userID int, orderIDs []int, action string) (int, error) {
	if action != "give" && action != "return" {
		return 0, ErrWrongArgument
	}

	ordersFailed := 0
	for i := range orderIDs {
		strOrderID := strconv.Itoa(orderIDs[i])
		someOrder, orderExists := s.data[strOrderID]
		if !orderExists {
			ordersFailed++
			continue
		}
		if someOrder.UserID != userID || !isOrderEligible(someOrder, action) {
			ordersFailed++
			continue
		}

		switch action {
		case "return":
			someOrder.Status = "returned"
		case "give":
			someOrder.Status = "given"
		}
		someOrder.LastChange = time.Now().Format(dateLayout)

		s.data[strOrderID] = someOrder
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
		if s.data[i].Status != "returned" {
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
