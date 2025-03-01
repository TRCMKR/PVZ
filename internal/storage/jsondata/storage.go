package jsondata

import (
	"errors"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/order"

	"github.com/bytedance/sonic"
)

var (
	errDataNotMarshalled = errors.New("data wasn't marshalled")
	errDataNotWritten    = errors.New("data wasn't written")
)

const (
	dateLayout = "2006.01.02 15:04:05"
)

const (
	orderGiven    = "given"
	orderReturned = "returned"
	orderStored   = "stored"
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

func (s *Storage) AddOrder(order order.Order) {
	stringOrderID := strconv.Itoa(order.ID)

	s.data[stringOrderID] = order
}

func (s *Storage) RemoveOrder(orderID int) {
	strOrderID := strconv.Itoa(orderID)

	delete(s.data, strOrderID)
}

func (s *Storage) UpdateOrder(orderID int, someOrder order.Order) {
	strOrderID := strconv.Itoa(orderID)

	s.data[strOrderID] = someOrder
}

func (s *Storage) GetByID(orderID int) order.Order {
	strOrderID := strconv.Itoa(orderID)

	return s.data[strOrderID]
}

func (s *Storage) GetByUserID(userID int) []order.Order {
	orderHistory := s.OrderHistory()
	userOrders := make([]order.Order, 0, len(orderHistory))
	for i := range orderHistory {
		if userID != orderHistory[i].UserID {
			continue
		}
		userOrders = append(userOrders, orderHistory[i])
	}

	return userOrders
}

func (s *Storage) GetReturns() []order.Order {
	returns := make([]order.Order, 0, len(s.data))

	for i := range s.data {
		if s.data[i].Status != orderReturned {
			continue
		}
		returns = append(returns, s.data[i])
	}

	return returns
}

func (s *Storage) OrderHistory() []order.Order {
	orderHistory := make([]order.Order, 0, len(s.data))
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

func (s *Storage) Contains(orderID int) bool {
	strOrderID := strconv.Itoa(orderID)
	_, ok := s.data[strOrderID]

	return ok
}
