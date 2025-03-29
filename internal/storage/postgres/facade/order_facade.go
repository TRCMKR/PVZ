package facade

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/cache/lru"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
	order_handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/order"
)

type orderStorage interface {
	AddOrder(context.Context, models.Order) error
	RemoveOrder(context.Context, int) error
	UpdateOrder(context.Context, int, models.Order) error
	GetByID(context.Context, int) (models.Order, error)
	GetByUserID(context.Context, int, int) ([]models.Order, error)
	GetReturns(context.Context) ([]models.Order, error)
	GetOrders(context.Context, []query.Cond, int, int) ([]models.Order, error)
	OffsetGetOrders(context.Context, []query.Cond, int, int, int) ([]models.Order, error)
	Contains(context.Context, int) (bool, error)
}

var (
	errWrongOperator = errors.New("wrong operator")
)

type OrderFacade struct {
	cache              *lru.Cache[int, models.Order]
	historyOrdersCache *lru.Cache[int, models.Order]
	orderStorage       orderStorage
}

func NewOrderFacade(ctx context.Context, orderStorage orderStorage, capacity int) *OrderFacade {
	historyOrdersCache := lru.NewCache[int, models.Order](capacity)

	recentOrders, _ := orderStorage.GetOrders(ctx, nil, capacity, 0)

	for _, order := range recentOrders {
		historyOrdersCache.Put(order.ID, order)
	}

	return &OrderFacade{
		orderStorage:       orderStorage,
		cache:              lru.NewCache[int, models.Order](capacity),
		historyOrdersCache: historyOrdersCache,
	}
}

func (f *OrderFacade) AddOrder(ctx context.Context, order models.Order) error {
	err := f.orderStorage.AddOrder(ctx, order)
	if err != nil {
		return err
	}

	f.cache.Put(order.ID, order)
	f.historyOrdersCache.Put(order.ID, order)

	return nil
}

func (f *OrderFacade) RemoveOrder(ctx context.Context, id int) error {
	err := f.orderStorage.RemoveOrder(ctx, id)
	if err != nil {
		return err
	}

	f.cache.Remove(id)
	f.historyOrdersCache.Remove(id)

	return nil
}

func (f *OrderFacade) UpdateOrder(ctx context.Context, id int, order models.Order) error {
	err := f.orderStorage.UpdateOrder(ctx, id, order)
	if err != nil {
		return err
	}

	f.cache.Put(order.ID, order)
	f.historyOrdersCache.Put(id, order)

	return nil
}

func (f *OrderFacade) GetByID(ctx context.Context, id int) (models.Order, error) {
	if order, ok := f.cache.Get(id); ok {
		return order, nil
	}

	order, err := f.orderStorage.GetByID(ctx, id)
	if err != nil {
		return models.Order{}, err
	}

	f.cache.Put(order.ID, order)

	return order, nil
}

func (f *OrderFacade) GetByUserID(ctx context.Context, id int, userID int) ([]models.Order, error) {
	return f.orderStorage.GetByUserID(ctx, id, userID)
}

func (f *OrderFacade) GetReturns(ctx context.Context) ([]models.Order, error) {
	return f.orderStorage.GetReturns(ctx)
}

func (f *OrderFacade) getOrderValue(field string, order models.Order) (interface{}, error) {
	switch field {
	case "id":
		return order.ID, nil
	case "user_id":
		return order.UserID, nil
	case "weight":
		return order.Weight, nil
	case "price":
		return order.Price, nil
	case "status":
		return order.Status, nil
	case "arrival_date":
		return order.ArrivalDate, nil
	case "expiry_date":
		return order.ExpiryDate, nil
	default:
		return nil, fmt.Errorf("unknown field %s", field)
	}
}

func (f *OrderFacade) checkOrder(order models.Order, cond query.Cond) (bool, error) {
	value, err := f.getOrderValue(cond.Field, order)
	if err != nil {
		return false, err
	}

	valueType := order_handler.GetFilters()[cond.Field]

	switch cond.Operator {
	case query.Equals:
		return value == cond.Value, nil
	case query.GreaterEqualThan:
		switch valueType {
		case order_handler.NumberType:
			return value.(int) >= cond.Value.(int), nil
		case order_handler.DateType:
			return value.(time.Time).After(cond.Value.(time.Time)) || value.(time.Time).Equal(cond.Value.(time.Time)), nil
		case order_handler.WordType:
			return value.(string) >= cond.Value.(string), nil
		}
	case query.GreaterThan:
		switch valueType {
		case order_handler.NumberType:
			return value.(int) > cond.Value.(int), nil
		case order_handler.DateType:
			return value.(time.Time).After(cond.Value.(time.Time)), nil
		case order_handler.WordType:
			return value.(string) > cond.Value.(string), nil
		}
	case query.LessEqualThan:
		switch valueType {
		case order_handler.NumberType:
			return value.(int) <= cond.Value.(int), nil
		case order_handler.DateType:
			return value.(time.Time).Before(cond.Value.(time.Time)) || value.(time.Time).Equal(cond.Value.(time.Time)), nil
		case order_handler.WordType:
			return value.(string) <= cond.Value.(string), nil
		}
	case query.LessThan:
		switch valueType {
		case order_handler.NumberType:
			return value.(int) < cond.Value.(int), nil
		case order_handler.DateType:
			return value.(time.Time).Before(cond.Value.(time.Time)), nil
		case order_handler.WordType:
			return value.(string) < cond.Value.(string), nil
		}
	case query.NotEquals:
		return value != cond.Value, nil
	}

	return false, errWrongOperator
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func (f *OrderFacade) GetOrders(ctx context.Context, params []query.Cond, count int, page int) ([]models.Order, error) {
	result := make([]models.Order, 0, count)

	operator := func(order models.Order) (bool, error) {
		var ok bool
		for _, cond := range params {
			var err error
			ok, err = f.checkOrder(order, cond)
			if err != nil {
				return false, err
			}

			if !ok {
				return false, nil
			}
		}

		return true, nil
	}

	recentOrders := f.historyOrdersCache.GetAllBy(operator)
	sort.Slice(recentOrders, func(i, j int) bool {
		return recentOrders[i].LastChange.After(recentOrders[j].LastChange)
	})

	params = append(params, query.Cond{
		Operator: query.GreaterThan,
		Field:    "last_change",
		Value:    recentOrders[len(recentOrders)-1].LastChange,
	})

	if count == 0 {
		dbOrders, err := f.orderStorage.GetOrders(ctx, params, count, 0)
		if err != nil {
			return nil, err
		}

		result = append(result, dbOrders...)

		return result, nil
	}

	newResult := make([]models.Order, 0, count)
	if len(result) < count && len(result) != 0 {
		dbOrders, err := f.orderStorage.GetOrders(ctx, params, count-len(result), 0)
		if err != nil {
			return nil, err
		}

		newResult = append(newResult, dbOrders...)
	} else if len(result) == 0 {
		//dbPage := len(result)/count + min(1, len(result)%count)
		dbPage := (len(result) + count - 1) / count

		dbOrders, err := f.orderStorage.OffsetGetOrders(ctx, params, count, dbPage, abs(len(result)-count*page))
		if err != nil {
			return nil, err
		}

		newResult = append(newResult, dbOrders...)
	} else {
		newResult = result[page*count : min(count*(page+1), len(result))]
	}

	return newResult, nil
}

func (f *OrderFacade) Contains(ctx context.Context, id int) (bool, error) {
	if _, ok := f.cache.Get(id); ok {
		return true, nil
	}
	if _, ok := f.historyOrdersCache.Get(id); ok {
		return true, nil
	}

	ok, err := f.orderStorage.Contains(ctx, id)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}

	order, _ := f.orderStorage.GetByID(ctx, id)

	f.cache.Put(order.ID, order)

	return true, nil
}
