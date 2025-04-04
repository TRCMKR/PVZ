package admin

import (
	"context"
	"errors"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

type adminStorage interface {
	CreateAdmin(context.Context, models.Admin) error
	GetAdminByUsername(context.Context, string) (models.Admin, error)
	UpdateAdmin(context.Context, int, models.Admin) error
	DeleteAdmin(context.Context, string) error
	ContainsUsername(context.Context, string) (bool, error)
	ContainsID(context.Context, int) (bool, error)
}

// Service ...
type Service struct {
	Storage adminStorage
}

var (
	// ErrUsernameUsed ...
	ErrUsernameUsed = errors.New("username is already used")

	// ErrIDUsed ...
	ErrIDUsed = errors.New("id is already used")

	// ErrAdminDoesntExist ...
	ErrAdminDoesntExist = errors.New("admin with such username doesn't exist")

	// ErrWrongPassword ...
	ErrWrongPassword = errors.New("wrong password")
)

// NewService ...
func NewService(storage adminStorage) *Service {
	return &Service{
		Storage: storage,
	}
}
