package admin

import (
	"context"
	"errors"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

type Handler struct {
	adminService adminService
}

func NewHandler(adminService adminService) *Handler {
	return &Handler{
		adminService: adminService,
	}
}

const (
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
	ErrFieldsMissing = errors.New("missing fields")
	ErrNoUsername    = errors.New("username wasn't provided")
)
