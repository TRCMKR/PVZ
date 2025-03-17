package order

import (
	"context"
	"errors"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
)

const (
	giveOrder   = "give"
	returnOrder = "return"
)

var (
	ErrOrderAlreadyExists = errors.New("such order exists")
	ErrOrderExpired       = errors.New("expired order")
	ErrOrderIsNotExpired  = errors.New("order is not expired")
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderIsGiven       = errors.New("order is given")
	ErrWrongWeight        = errors.New("wrong weight")
	ErrWrongPrice         = errors.New("wrong price")
	ErrOrderNotEligible   = errors.New("order not eligible")
	ErrUndefinedAction    = errors.New("undefined action")
	ErrNotEnoughWeight    = errors.New("not enough weight")
	ErrWrongPackaging     = errors.New("wrong packaging")
	ErrUndefinedPackaging = errors.New("undefined packaging")
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

type Service struct {
	Storage orderStorage
}

func NewService(storage orderStorage) *Service {
	return &Service{
		Storage: storage,
	}
}
