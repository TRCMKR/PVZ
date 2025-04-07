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

// Service is a struct for admin service
type Service struct {
	Storage adminStorage
}

var (
	// ErrUsernameUsed happens when username is used
	ErrUsernameUsed = errors.New("username is already used")

	// ErrIDUsed happens when id is used
	ErrIDUsed = errors.New("id is already used")

	// ErrAdminDoesntExist happens when admin doesn't exist
	ErrAdminDoesntExist = errors.New("admin with such username doesn't exist")

	// ErrWrongPassword happens when wrong password was passed
	ErrWrongPassword = errors.New("wrong password")
)

// NewService creates instance of admin Service
func NewService(storage adminStorage) *Service {
	return &Service{
		Storage: storage,
	}
}
