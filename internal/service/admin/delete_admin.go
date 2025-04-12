package admin

import (
	"context"

	"go.uber.org/zap"
)

// DeleteAdmin deletes admin
func (s *Service) DeleteAdmin(ctx context.Context, password string, username string) error {
	ok, err := s.ContainsUsername(ctx, username)
	if err != nil {
		return err
	}
	if !ok {
		s.logger.Error(ErrAdminDoesntExist.Error(),
			zap.String("username", username),
		)

		return ErrAdminDoesntExist
	}

	admin, err := s.Storage.GetAdminByUsername(ctx, username)
	if err != nil {
		return err
	}

	if !admin.CheckPassword(password) {
		s.logger.Error(ErrWrongPassword.Error(),
			zap.String("username", username),
			zap.Error(ErrWrongPassword),
		)

		return ErrWrongPassword
	}

	return s.Storage.DeleteAdmin(ctx, username)
}
