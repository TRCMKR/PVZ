package order

import (
	"context"
	"errors"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	myquery "gitlab.ozon.dev/alexplay1224/homework/internal/query"

	"github.com/Rhymond/go-money"
)

// Handler ...
type Handler struct {
	OrderService orderService
}

// NewHandler ...
func NewHandler(orderService orderService) *Handler {
	return &Handler{
		OrderService: orderService,
	}
}

const (
	// OrderIDParam ...
	OrderIDParam = "id"

	// UserIDParam ...
	UserIDParam = "user_id"

	// WeightParam ...
	WeightParam = "weight"

	// PriceParam ...
	PriceParam = "price"

	// StatusParam ...
	StatusParam = "status"

	// ArrivalDateParam ...
	ArrivalDateParam = "arrival_date"

	// ArrivalDateFromParam ...
	ArrivalDateFromParam = "arrival_date_from"

	// ArrivalDateToParam ...
	ArrivalDateToParam = "arrival_date_to"

	// ExpiryDateParam ...
	ExpiryDateParam = "expiry_date"

	// ExpiryDateFromParam ...
	ExpiryDateFromParam = "expiry_date_from"

	// ExpiryDateToParam ...
	ExpiryDateToParam = "expiry_date_to"

	// WeightFromParam ...
	WeightFromParam = "weight_from"

	// WeightToParam ...
	WeightToParam = "weight_to"

	// PriceFromParam ...
	PriceFromParam = "price_from"

	// PriceToParam ...
	PriceToParam = "price_to"

	// CountParam ...
	CountParam = "count"

	// PageParam ...
	PageParam = "page"
)

type orderService interface {
	AcceptOrder(context.Context, int, int, float64, money.Money, time.Time, []models.Packaging) error
	ReturnOrder(context.Context, int) error
	ProcessOrder(context.Context, int, int, string) error
	UserOrders(context.Context, int, int) ([]models.Order, error)
	Returns(context.Context) ([]models.Order, error)
	GetOrders(context.Context, []myquery.Cond, int, int) ([]models.Order, error)
}

var (
	errNoSuchPackaging   = errors.New("no such packaging")
	errInvalidOrderID    = errors.New("invalid order id")
	errWrongNumberFormat = errors.New("wrong number format")
	errWrongDateFormat   = errors.New("wrong date format")
	errWrongStatusFormat = errors.New("wrong status format")
	errFieldsMissing     = errors.New("missing fields")
	errWrongJSONFormat   = errors.New("wrong json format")
)

// InputType ...
type InputType uint

const (
	// NumberType ...
	NumberType InputType = iota
	// WordType ...
	WordType
	// DateType ...
	DateType
)

const (
	inputDateAndTimeLayout = "2006.01.02-15:04:05"
	inputDateLayout        = "2006.01.02"
)
