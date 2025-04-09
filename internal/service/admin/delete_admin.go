package admin

import (
	"context"
)

// DeleteAdmin deletes admin
func (s *Service) DeleteAdmin(ctx context.Context, password string, username string) error {
	ok, err := s.ContainsUsername(ctx, username)
	if err != nil {
		return err
	}
	if !ok {
		return ErrAdminDoesntExist
	}

	admin, err := s.Storage.GetAdminByUsername(ctx, username)
	if err != nil {
		return err
	}

	if !admin.CheckPassword(password) {
		return ErrWrongPassword
	}

	return s.Storage.DeleteAdmin(ctx, username)
}
