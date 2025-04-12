package admin

import (
	"context"
	"errors"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// Handler is a struct for handling admin related call
type Handler struct {
	adminService adminService
}

// NewHandler creates an instance of admin Handler
func NewHandler(adminService adminService) *Handler {
	return &Handler{
		adminService: adminService,
	}
}

const (
	// AdminUsernameParam is a query param for admin username
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
	// ErrFieldsMissing happens when some fields are missing
	ErrFieldsMissing = errors.New("missing fields")

	// ErrNoUsername happens when username wasn't provided
	ErrNoUsername = errors.New("username wasn't provided")
)
