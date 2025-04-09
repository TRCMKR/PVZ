package order

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
)

const (
	giveOrder   = "give"
	returnOrder = "return"
)

var (
	// ErrOrderAlreadyExists happens when such order exists
	ErrOrderAlreadyExists = errors.New("such order exists")

	// ErrOrderExpired happens when order expired
	ErrOrderExpired = errors.New("expired order")

	// ErrOrderIsNotExpired happens when order is not yet expired
	ErrOrderIsNotExpired = errors.New("order is not expired")

	// ErrOrderNotFound happens when such order is not found
	ErrOrderNotFound = errors.New("order not found")

	// ErrOrderIsGiven happens when order is given
	ErrOrderIsGiven = errors.New("order is given")

	// ErrWrongWeight happens when order weight doesn't satisfy
	ErrWrongWeight = errors.New("wrong weight")

	// ErrWrongPrice happens when order has incorrect price
	ErrWrongPrice = errors.New("wrong price")

	// ErrOrderNotEligible happens when order is not eligible
	ErrOrderNotEligible = errors.New("order not eligible")

	// ErrUndefinedAction happens when undefined action is passed
	ErrUndefinedAction = errors.New("undefined action")

	// ErrNotEnoughWeight happens when order doesn't have enough weight
	ErrNotEnoughWeight = errors.New("not enough weight")

	// ErrWrongPackaging happens when packaging is wrong
	ErrWrongPackaging = errors.New("wrong packaging")
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
}

type txManager interface {
	RunSerializable(context.Context, func(context.Context, pgx.Tx) error) error
	RunRepeatableRead(context.Context, func(context.Context, pgx.Tx) error) error
	RunReadCommitted(context.Context, func(context.Context, pgx.Tx) error) error
}

// Service is a structure for order service
type Service struct {
	Storage   orderStorage
	txManager txManager
}

// NewService creates instance of an order Service
func NewService(storage orderStorage, txManager txManager) *Service {
	return &Service{
		Storage:   storage,
		txManager: txManager,
	}
}
