package jsondata

import (
	"context"
	"errors"
	"log"
	"os"
	"sort"
	"strconv"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/bytedance/sonic"
)

var (
	errDataNotMarshalled = errors.New("data wasn't marshalled")
	errDataNotWritten    = errors.New("data wasn't written")
)

const (
	orderReturned = "returned"
)

type Storage struct {
	data map[string]models.Order
	path string
}

func New(path string) (*Storage, error) {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("read json fail", err)
	}

	data := make(map[string]models.Order)

	err = sonic.Unmarshal(jsonData, &data)
	if err != nil {
		log.Fatal("unmarshal json fail", err)
	}

	return &Storage{
		data: data,
		path: path,
	}, nil
}

func (s *Storage) AddOrder(_ context.Context, order models.Order) {
	stringOrderID := strconv.Itoa(order.ID)

	s.data[stringOrderID] = order
}

func (s *Storage) RemoveOrder(_ context.Context, orderID int) {
	strOrderID := strconv.Itoa(orderID)

	delete(s.data, strOrderID)
}

func (s *Storage) UpdateOrder(_ context.Context, orderID int, someOrder models.Order) {
	strOrderID := strconv.Itoa(orderID)

	s.data[strOrderID] = someOrder
}

func (s *Storage) GetByID(_ context.Context, orderID int) models.Order {
	strOrderID := strconv.Itoa(orderID)

	return s.data[strOrderID]
}

func (s *Storage) GetByUserID(ctx context.Context, userID int, count int) []models.Order {
	orderHistory := s.GetOrders(ctx, nil, 0, 0)
	userOrders := make([]models.Order, 0, len(orderHistory))
	currentCount := 0
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

	return userOrders
}

func (s *Storage) GetReturns(_ context.Context) []models.Order {
	returns := make([]models.Order, 0, len(s.data))

	for i := range s.data {
		if s.data[i].Status != orderReturned {
			continue
		}
		returns = append(returns, s.data[i])
	}

	return returns
}

func (s *Storage) GetOrders(_ context.Context, _ map[string]string, _ int, _ int) []models.Order {
	orderHistory := make([]models.Order, 0, len(s.data))
	for i := range s.data {
		orderHistory = append(orderHistory, s.data[i])
	}
	sort.Slice(orderHistory, func(i, j int) bool {
		return orderHistory[j].LastChange.Before(orderHistory[i].LastChange)
	})

	return orderHistory
}

func (s *Storage) Save(_ context.Context) error {
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

func (s *Storage) Contains(_ context.Context, orderID int) bool {
	strOrderID := strconv.Itoa(orderID)
	_, ok := s.data[strOrderID]

	return ok
}
