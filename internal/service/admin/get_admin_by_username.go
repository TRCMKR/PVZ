package admin

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (s *Service) GetAdminByUsername(ctx context.Context, username string) (models.Admin, error) {
	ok, err := s.ContainsUsername(ctx, username)
	if err != nil {
		return models.Admin{}, err
	}
	if !ok {
		return models.Admin{}, ErrAdminDoesntExist
	}

	return s.Storage.GetAdminByUsername(ctx, username)
}
