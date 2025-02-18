package json_data

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/bytedance/sonic"

	"homework/order"
)

var (
	ErrorOrderAlreadyExists  = errors.New("such order exists")
	ErrorOrderExpired        = errors.New("expired order")
	ErrorWrongDateFormat     = errors.New("wrong date format")
	ErrorDataNotMarshalled   = errors.New("data wasn't marshalled")
	ErrorDataNotUnmarshalled = errors.New("data wasn't unmarshalled")
	ErrorDataNotWritten      = errors.New("data wasn't written")
	ErrorOrderNotFound       = errors.New("order not found")
	ErrorWrongArgument       = errors.New("wrong argument")
	ErrorFileNotOpened       = errors.New("file not opened")
)

const (
	dateLayout       = "2006.01.02 15:04:05"
	inputDateLayout1 = "2006.01.02-15:04:05"
	inputDateLayout2 = "2006.01.02"
)

type Storage struct {
	data []*order.Order
	path string
}

func New(path string) (*Storage, error) {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("read json fail", err)
		os.Exit(1)
	}

	data := make([]*order.Order, 0)

	err = sonic.Unmarshal(jsonData, &data)
	if err != nil {
		fmt.Println("unmarshal json fail", err)
		os.Exit(1)
	}

	return &Storage{
		data: data,
		path: path,
	}, nil
}

func (s *Storage) Save() error {
	jsonData, err := sonic.MarshalIndent(s.data, "", "  ")
	if err != nil {
		return ErrorDataNotMarshalled
	}
	err = os.WriteFile(s.path, jsonData, 0644)
	if err != nil {
		return ErrorDataNotWritten
	}

	return nil
}

func (s *Storage) AcceptOrder(orderID string, userID string, expiryDate string) error {
	date, err := time.Parse(inputDateLayout1, expiryDate)
	if err != nil {
		date, err = time.Parse(inputDateLayout2, expiryDate)
		if err != nil {
			return ErrorWrongDateFormat
		}
	}
	if date.Before(time.Now()) {
		return ErrorOrderExpired
	}

	for _, o := range s.data {
		if o.OrderID == orderID {
			return ErrorOrderAlreadyExists
		}
	}

	currentTime := time.Now().Format(dateLayout)

	s.data = append(s.data, &order.Order{
		OrderID:     orderID,
		UserId:      userID,
		ArrivalDate: currentTime,
		Status:      "stored",
		ExpiryDate:  date.Format(dateLayout),
		LastChange:  currentTime,
	})

	return nil
}

func (s *Storage) AcceptOrders(path string) (int, error) {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return 0, ErrorFileNotOpened
	}

	data := make([]*order.Order, 0)

	err = sonic.Unmarshal(jsonData, &data)
	if err != nil {
		return 0, ErrorDataNotUnmarshalled
	}

	orderMap := make(map[string]*order.Order)
	for _, o := range s.data {
		orderMap[o.OrderID] = o
	}

	ordersFailed := 0
	for _, o := range data {
		if _, ok := orderMap[o.OrderID]; !ok {
			s.data = append(s.data, o)
			continue
		}

		ordersFailed++
	}

	return ordersFailed, nil
}

func (s *Storage) ReturnOrder(orderID string) error {
	ind := -1
	for i, v := range s.data {
		if v.OrderID == orderID {
			ind = i
			break
		}
	}
	if ind == -1 {
		return ErrorOrderNotFound
	}

	s.data[ind] = s.data[len(s.data)-1]
	s.data = s.data[:len(s.data)-1]

	return nil
}

func (s *Storage) ProcessOrders(userID string, orderIDs []string, action string) (int, error) {
	if action != "give" && action != "return" {
		return 0, ErrorWrongArgument
	}

	userOrders, _ := s.UserOrders(userID)
	orderMap := make(map[string]*order.Order)
	for _, o := range userOrders {
		orderMap[o.OrderID] = o
	}

	ordersFailed := 0
	deadline := time.Now()
	if action == "return" {
		deadline = deadline.AddDate(0, 0, -2)
	}
	for _, o := range orderIDs {
		if someOrder, ok := orderMap[o]; ok {
			var date time.Time
			if action == "return" {
				date, _ = time.Parse(dateLayout, someOrder.LastChange)
			} else {
				date, _ = time.Parse(dateLayout, someOrder.ExpiryDate)
			}

			if action == "give" && someOrder.Status == "stored" && date.After(deadline) ||
				action == "return" && someOrder.Status == "given" && date.After(deadline) {
				if action == "give" {
					someOrder.Status = "given"
				} else {
					someOrder.Status = "returned"
				}
				someOrder.LastChange = time.Now().Format(dateLayout)
				continue
			}
		}

		ordersFailed++
	}

	return ordersFailed, nil
}

func (s *Storage) UserOrders(args ...string) ([]*order.Order, error) {
	userID := args[0]
	var n int
	var err error
	if n, err = strconv.Atoi(args[len(args)-1]); len(args) == 2 && err != nil {
		return nil, ErrorWrongArgument
	}
	userOrders := make([]*order.Order, 0)

	count := 0
	orderHistory := s.OrderHistory()
	for _, v := range orderHistory {
		if userID != v.UserId {
			continue
		}
		userOrders = append(userOrders, v)
		count++
		if count == n {
			break
		}
	}

	return userOrders, nil
}

func (s *Storage) Returns() []*order.Order {
	returns := make([]*order.Order, 0)

	for _, v := range s.data {
		if v.Status != "returned" {
			continue
		}
		returns = append(returns, v)
	}

	return returns
}

func (s *Storage) OrderHistory() []*order.Order {
	orderHistory := s.data[:]
	sort.Slice(orderHistory, func(i, j int) bool {
		t1, _ := time.Parse(dateLayout, orderHistory[i].LastChange)
		t2, _ := time.Parse(dateLayout, orderHistory[j].LastChange)
		return t2.Before(t1)
	})

	return orderHistory
}
