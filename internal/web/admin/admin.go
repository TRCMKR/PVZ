package admin

import (
	"context"
	"errors"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// Handler ...
type Handler struct {
	adminService adminService
}

// NewHandler ...
func NewHandler(adminService adminService) *Handler {
	return &Handler{
		adminService: adminService,
	}
}

const (
	// AdminUsernameParam ...
	AdminUsernameParam = "admin_username"
)

type adminService interface {
	CreateAdmin(context.Context, models.Admin) error
	GetAdminByUsername(context.Context, string) (models.Admin, error)
	UpdateAdmin(context.Context, string, string, models.Admin) error
	DeleteAdmin(context.Context, string, string) error
	ContainsUsername(context.Context, string) (bool, error)
	ContainsID(context.Context, int) (bool, error)
}

var (
	// ErrFieldsMissing ...
	ErrFieldsMissing = errors.New("missing fields")

	// ErrNoUsername ...
	ErrNoUsername = errors.New("username wasn't provided")
)
