package admin

import (
	"context"

	"go.uber.org/zap"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// UpdateAdmin updates admin with new fields
func (s *Service) UpdateAdmin(ctx context.Context, username string, password string, admin models.Admin) error {
	ok, err := s.ContainsUsername(ctx, username)
	if err != nil {
		return err
	}
	if !ok {
		s.logger.Error(ErrAdminDoesntExist.Error(),
			zap.String("username", username),
			zap.Error(ErrAdminDoesntExist),
		)

		return ErrAdminDoesntExist
	}

	someAdmin, err := s.GetAdminByUsername(ctx, username)
	if err != nil {
		return err
	}
	if !someAdmin.CheckPassword(password) {
		s.logger.Error(ErrWrongPassword.Error(),
			zap.String("username", username),
			zap.Error(ErrWrongPassword),
		)

		return ErrWrongPassword
	}

	return s.UpdateAdmin(ctx, username, password, admin)
}
