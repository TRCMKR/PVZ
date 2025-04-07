package admin

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// UpdateAdmin updates admin with new fields
func (s *Service) UpdateAdmin(ctx context.Context, username string, password string, admin models.Admin) error {
	ok, err := s.ContainsUsername(ctx, username)
	if err != nil {
		return err
	}
	if !ok {
		return ErrAdminDoesntExist
	}

	someAdmin, err := s.GetAdminByUsername(ctx, username)
	if err != nil {
		return err
	}
	if !someAdmin.CheckPassword(password) {
		return ErrWrongPassword
	}

	return s.UpdateAdmin(ctx, username, password, admin)
}
