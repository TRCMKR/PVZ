//go:generate mockgen -source=main.go -destination=../../mocks/service/mock_order_service.go -package=service

package order

import (
	"context"
	"errors"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	myquery "gitlab.ozon.dev/alexplay1224/homework/internal/query"

	"github.com/Rhymond/go-money"
)

type Handler struct {
	OrderService orderService
}

func NewHandler(orderService orderService) *Handler {
	return &Handler{
		OrderService: orderService,
	}
}

const (
	OrderIDParam         = "id"
	UserIDParam          = "user_id"
	WeightParam          = "weight"
	PriceParam           = "price"
	StatusParam          = "status"
	ArrivalDateParam     = "arrival_date"
	ArrivalDateFromParam = "arrival_date_from"
	ArrivalDateToParam   = "arrival_date_to"
	ExpiryDateParam      = "expiry_date"
	ExpiryDateFromParam  = "expiry_date_from"
	ExpiryDateToParam    = "expiry_date_to"
	WeightFromParam      = "weight_from"
	WeightToParam        = "weight_to"
	PriceFromParam       = "price_from"
	PriceToParam         = "price_to"
	CountParam           = "count"
	PageParam            = "page"
)

type orderService interface {
	AcceptOrder(context.Context, int, int, float64, money.Money, time.Time, []models.Packaging) error
	AcceptOrders(context.Context, map[string]models.Order) (int, error)
	ReturnOrder(context.Context, int) error
	ProcessOrders(context.Context, int, []int, string) (int, error)
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
	errWrongJsonFormat   = errors.New("wrong json format")
)

type inputType uint

const (
	numberType inputType = iota
	wordType
	dateType
)

const (
	inputDateAndTimeLayout = "2006.01.02-15:04:05"
	inputDateLayout        = "2006.01.02"
)
