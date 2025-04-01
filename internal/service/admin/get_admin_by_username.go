package admin

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// GetAdminByUsername ...
func (s *Service) GetAdminByUsername(ctx context.Context, username string) (models.Admin, error) {
	var admin models.Admin

	ok, err := s.ContainsUsername(ctx, username)
	if err != nil {
		return models.Admin{}, err
	}
	if !ok {
		return models.Admin{}, ErrAdminDoesntExist
	}

	admin, err = s.Storage.GetAdminByUsername(ctx, username)
	if err != nil {
		return models.Admin{}, err
	}

	return admin, nil
}
