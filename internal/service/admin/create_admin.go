package admin

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// CreateAdmin ...
func (s *Service) CreateAdmin(ctx context.Context, admin models.Admin) error {
	ok, err := s.ContainsUsername(ctx, admin.Username)
	if err != nil {
		return err
	}
	if ok {
		return ErrUsernameUsed
	}

	ok, err = s.ContainsID(ctx, admin.ID)
	if err != nil {
		return err
	}
	if ok {
		return ErrIDUsed
	}

	return s.Storage.CreateAdmin(ctx, admin)
}
