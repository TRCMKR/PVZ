package admin

import (
	"context"

	"go.uber.org/zap"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// GetAdminByUsername gets admin by username
func (s *Service) GetAdminByUsername(ctx context.Context, username string) (models.Admin, error) {
	var admin models.Admin

	ok, err := s.ContainsUsername(ctx, username)
	if err != nil {
		return models.Admin{}, err
	}
	if !ok {
		s.logger.Error(ErrAdminDoesntExist.Error(),
			zap.String("username", username),
			zap.Error(ErrAdminDoesntExist),
		)

		return models.Admin{}, ErrAdminDoesntExist
	}

	admin, err = s.Storage.GetAdminByUsername(ctx, username)
	if err != nil {
		return models.Admin{}, err
	}

	return admin, nil
}
