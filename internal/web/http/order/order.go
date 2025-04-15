package order

import (
	"context"
	"errors"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	myquery "gitlab.ozon.dev/alexplay1224/homework/internal/query"

	"github.com/Rhymond/go-money"
)

// Handler is a structure for order handling
type Handler struct {
	OrderService orderService
}

// NewHandler creates an instance of order Handler
func NewHandler(orderService orderService) *Handler {
	return &Handler{
		OrderService: orderService,
	}
}

const (
	// OrderIDParam is a param for order id
	OrderIDParam = "id"

	// UserIDParam is a param for user id
	UserIDParam = "user_id"

	// WeightParam is a param for weight
	WeightParam = "weight"

	// PriceParam is a param for price
	PriceParam = "price"

	// StatusParam is a param for status
	StatusParam = "status"

	// ArrivalDateParam is a param for arrival date
	ArrivalDateParam = "arrival_date"

	// ArrivalDateFromParam is a param for arrival date from
	ArrivalDateFromParam = "arrival_date_from"

	// ArrivalDateToParam is a param for arrival date to
	ArrivalDateToParam = "arrival_date_to"

	// ExpiryDateParam is a param for expiry date
	ExpiryDateParam = "expiry_date"

	// ExpiryDateFromParam is a param for expiry date from
	ExpiryDateFromParam = "expiry_date_from"

	// ExpiryDateToParam is a param for expiry date to
	ExpiryDateToParam = "expiry_date_to"

	// WeightFromParam is a param for weight from
	WeightFromParam = "weight_from"

	// WeightToParam is a param for weight to
	WeightToParam = "weight_to"

	// PriceFromParam is a param for price from
	PriceFromParam = "price_from"

	// PriceToParam is a param for price to
	PriceToParam = "price_to"

	// CountParam is a param for count
	CountParam = "count"

	// PageParam is a param for page
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

// InputType is a type for all inputs
type InputType uint

const (
	// NumberType is for number inputs
	NumberType InputType = iota
	// WordType is for word inputs
	WordType
	// DateType is for date inputs
	DateType
)

const (
	inputDateAndTimeLayout = "2006.01.02-15:04:05"
	inputDateLayout        = "2006.01.02"
)
