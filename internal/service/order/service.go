package order

import (
	"context"
	"errors"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/tx_manager"
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
	AddOrder(context.Context, models.Order) error
	RemoveOrder(context.Context, int) error
	UpdateOrder(context.Context, int, models.Order) error
	GetByID(context.Context, int) (models.Order, error)
	GetByUserID(context.Context, int, int) ([]models.Order, error)
	GetReturns(context.Context) ([]models.Order, error)
	GetOrders(context.Context, []query.Cond, int, int) ([]models.Order, error)
	Contains(context.Context, int) (bool, error)
}

type txManager interface {
	RunSerializable(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunRepeatableRead(ctx context.Context, fn func(ctxTx context.Context) error) error
	RunReadCommitted(ctx context.Context, fn func(ctxTx context.Context) error) error
	GetQueryEngine(ctx context.Context) tx_manager.Database
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
