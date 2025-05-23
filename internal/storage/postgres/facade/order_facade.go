package facade

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/opentracing/opentracing-go"

	"gitlab.ozon.dev/alexplay1224/homework/internal/cache/lru"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
)

type orderStorage interface {
	AddOrder(context.Context, pgx.Tx, models.Order) error
	RemoveOrder(context.Context, pgx.Tx, int) error
	UpdateOrder(context.Context, pgx.Tx, int, models.Order) error
	GetByID(context.Context, pgx.Tx, int) (models.Order, error)
	GetByUserID(context.Context, pgx.Tx, int, int) ([]models.Order, error)
	GetReturns(context.Context, pgx.Tx) ([]models.Order, error)
	GetOrders(context.Context, pgx.Tx, []query.Cond, int, int) ([]models.Order, error)
	Contains(context.Context, pgx.Tx, int) (bool, error)
	OffsetGetOrders(context.Context, pgx.Tx, []query.Cond, int, int, int) ([]models.Order, error)
}

var (
	errWrongOperator        = errors.New("wrong operator")
	errUnsupportedValueType = errors.New("unsupported value type")
)

// OrderFacade is a structure for order facade
type OrderFacade struct {
	cache              *lru.Cache[int, models.Order]
	historyOrdersCache *lru.Cache[int, models.Order]
	orderStorage       orderStorage
}

// NewOrderFacade creates instance for order facade
func NewOrderFacade(ctx context.Context, orderStorage orderStorage, capacity int) *OrderFacade {
	historyOrdersCache := lru.NewCache[int, models.Order](capacity)

	recentOrders, _ := orderStorage.GetOrders(ctx, nil, nil, capacity, 0)

	for _, order := range recentOrders {
		historyOrdersCache.Put(order.ID, order)
	}

	return &OrderFacade{
		orderStorage:       orderStorage,
		cache:              lru.NewCache[int, models.Order](capacity),
		historyOrdersCache: historyOrdersCache,
	}
}

// AddOrder adds order
func (f *OrderFacade) AddOrder(ctx context.Context, tx pgx.Tx, order models.Order) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderFacade.AddOrder")
	defer span.Finish()

	err := f.orderStorage.AddOrder(ctx, tx, order)
	if err != nil {
		span.SetTag("error", err)

		return err
	}

	f.cache.Put(order.ID, order)
	f.historyOrdersCache.Put(order.ID, order)

	return nil
}

// RemoveOrder removes order
func (f *OrderFacade) RemoveOrder(ctx context.Context, tx pgx.Tx, id int) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderFacade.RemoveOrder")
	defer span.Finish()

	err := f.orderStorage.RemoveOrder(ctx, tx, id)
	if err != nil {
		span.SetTag("error", err)

		return err
	}

	f.cache.Remove(id)
	f.historyOrdersCache.Remove(id)

	return nil
}

// UpdateOrder updates order
func (f *OrderFacade) UpdateOrder(ctx context.Context, tx pgx.Tx, id int, order models.Order) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderFacade.UpdateOrder")
	defer span.Finish()

	err := f.orderStorage.UpdateOrder(ctx, tx, id, order)
	if err != nil {
		span.SetTag("error", err)

		return err
	}

	f.cache.Put(order.ID, order)
	f.historyOrdersCache.Put(id, order)

	return nil
}

// GetByID gets order by id
func (f *OrderFacade) GetByID(ctx context.Context, tx pgx.Tx, id int) (models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderFacade.GetByID")
	defer span.Finish()

	if order, ok := f.cache.Get(id); ok {
		span.SetTag("cache", true)

		return order, nil
	}

	order, err := f.orderStorage.GetByID(ctx, tx, id)
	if err != nil {
		span.SetTag("error", err)

		return models.Order{}, err
	}

	f.cache.Put(order.ID, order)

	return order, nil
}

// GetByUserID gets orders by user id
func (f *OrderFacade) GetByUserID(ctx context.Context, tx pgx.Tx, id int, userID int) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderFacade.GetByUserID")
	defer span.Finish()

	return f.orderStorage.GetByUserID(ctx, tx, id, userID)
}

// GetReturns gets all returned orders
func (f *OrderFacade) GetReturns(ctx context.Context, tx pgx.Tx) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderFacade.GetReturns")
	defer span.Finish()

	return f.orderStorage.GetReturns(ctx, tx)
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

	return f.compareValues(value, cond)
}

func (f *OrderFacade) compareValues(value interface{}, cond query.Cond) (bool, error) {
	switch cond.Operator {
	case query.Equals:
		return value == cond.Value, nil
	case query.NotEquals:
		return value != cond.Value, nil
	default:
		return f.compare(value, cond.Value, cond.Operator)
	}
}

func (f *OrderFacade) compare(value interface{}, condValue interface{}, operator query.CondType) (bool, error) {
	switch value := value.(type) {
	case int:
		condValueInt, _ := condValue.(int)

		return f.compareInts(value, condValueInt, operator)
	case time.Time:
		condValueTime, _ := condValue.(time.Time)

		return f.compareTimes(value, condValueTime, operator)
	case string:
		condValueStr, _ := condValue.(string)

		return f.compareStrings(value, condValueStr, operator)
	}

	return false, errUnsupportedValueType
}

func (f *OrderFacade) compareInts(value, condValue int, operator query.CondType) (bool, error) {
	switch operator {
	case query.GreaterEqualThan:
		return value >= condValue, nil
	case query.GreaterThan:
		return value > condValue, nil
	case query.LessEqualThan:
		return value <= condValue, nil
	case query.LessThan:
		return value < condValue, nil
	default:
		return false, errWrongOperator
	}
}

func (f *OrderFacade) compareTimes(value, condValue time.Time, operator query.CondType) (bool, error) {
	switch operator {
	case query.GreaterEqualThan:
		return value.After(condValue) || value.Equal(condValue), nil
	case query.GreaterThan:
		return value.After(condValue), nil
	case query.LessEqualThan:
		return value.Before(condValue) || value.Equal(condValue), nil
	case query.LessThan:
		return value.Before(condValue), nil
	default:
		return false, errWrongOperator
	}
}

func (f *OrderFacade) compareStrings(value, condValue string, operator query.CondType) (bool, error) {
	switch operator {
	case query.GreaterEqualThan:
		return value >= condValue, nil
	case query.GreaterThan:
		return value > condValue, nil
	case query.LessEqualThan:
		return value <= condValue, nil
	case query.LessThan:
		return value < condValue, nil
	default:
		return false, errWrongOperator
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func (f *OrderFacade) checkIfOrderMatches(order models.Order, params []query.Cond) (bool, error) {
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

// GetOrders gets orders by conditions
func (f *OrderFacade) GetOrders(ctx context.Context, tx pgx.Tx, params []query.Cond,
	count int, page int) ([]models.Order, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderFacade.GetOrders")
	defer span.Finish()

	recentOrders, err := f.historyOrdersCache.GetAllBy(func(order models.Order) (bool, error) {
		return f.checkIfOrderMatches(order, params)
	})
	if err != nil {
		span.SetTag("error", err)

		return nil, err
	}

	sort.Slice(recentOrders, func(i, j int) bool {
		return recentOrders[i].LastChange.After(recentOrders[j].LastChange)
	})

	result := make([]models.Order, 0, count)
	if count == 0 {
		dbOrders, err := f.orderStorage.GetOrders(ctx, tx, params, count, 0)
		if err != nil {
			span.SetTag("error", err)

			return nil, err
		}

		recentOrders = append(recentOrders, dbOrders...)

		span.SetTag("cache", true)

		return recentOrders, nil
	}

	recentOrders = recentOrders[min(page*count, len(recentOrders)):min(count*(page+1), len(recentOrders))]

	switch {
	case len(recentOrders) < count && len(recentOrders) != 0:
		params = append(params, query.Cond{
			Operator: query.LessThan,
			Field:    "last_change",
			Value:    recentOrders[len(recentOrders)-1].LastChange,
		})

		dbOrders, err := f.orderStorage.GetOrders(ctx, tx, params, count-len(result), 0)
		if err != nil {
			span.SetTag("error", err)

			return nil, err
		}
		recentOrders = append(recentOrders, dbOrders...)

	case len(recentOrders) == 0:
		dbPage := (len(result) + count - 1) / count
		dbOrders, err := f.orderStorage.OffsetGetOrders(ctx, tx, params, count, dbPage, abs(len(result)-count*page))
		if err != nil {
			span.SetTag("error", err)

			return nil, err
		}
		recentOrders = append(recentOrders, dbOrders...)

	default:
		result = recentOrders
		span.SetTag("cache", true)
	}

	return result, nil
}

// Contains checks if order is present
func (f *OrderFacade) Contains(ctx context.Context, tx pgx.Tx, id int) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "orderFacade.Contains")
	defer span.Finish()

	if _, ok := f.cache.Get(id); ok {
		span.SetTag("cache", true)

		return true, nil
	}
	if _, ok := f.historyOrdersCache.Get(id); ok {
		span.SetTag("cache", true)

		return true, nil
	}

	ok, err := f.orderStorage.Contains(ctx, tx, id)
	if err != nil {
		span.SetTag("error", err)

		return false, err
	}
	if !ok {
		return false, nil
	}

	order, _ := f.orderStorage.GetByID(ctx, tx, id)

	f.cache.Put(order.ID, order)

	return true, nil
}
