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
	// ErrOrderAlreadyExists ...
	ErrOrderAlreadyExists = errors.New("such order exists")

	// ErrOrderExpired ...
	ErrOrderExpired = errors.New("expired order")

	// ErrOrderIsNotExpired ...
	ErrOrderIsNotExpired = errors.New("order is not expired")

	// ErrOrderNotFound ...
	ErrOrderNotFound = errors.New("order not found")

	// ErrOrderIsGiven ...
	ErrOrderIsGiven = errors.New("order is given")

	// ErrWrongWeight ...
	ErrWrongWeight = errors.New("wrong weight")

	// ErrWrongPrice ...
	ErrWrongPrice = errors.New("wrong price")

	// ErrOrderNotEligible ...
	ErrOrderNotEligible = errors.New("order not eligible")

	// ErrUndefinedAction ...
	ErrUndefinedAction = errors.New("undefined action")

	// ErrNotEnoughWeight ...
	ErrNotEnoughWeight = errors.New("not enough weight")

	// ErrWrongPackaging ...
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

// Service ...
type Service struct {
	Storage   orderStorage
	txManager txManager
}

// NewService ...
func NewService(storage orderStorage, txManager txManager) *Service {
	return &Service{
		Storage:   storage,
		txManager: txManager,
	}
}
