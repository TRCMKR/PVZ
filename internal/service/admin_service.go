package service

import (
	"context"
	"errors"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

type adminStorage interface {
	CreateAdmin(context.Context, models.Admin) error
	GetAdminByUsername(context.Context, string) (models.Admin, error)
	UpdateAdmin(context.Context, int, models.Admin) error
	DeleteAdmin(context.Context, string) error
	ContainsUsername(context.Context, string) (bool, error)
	ContainsID(context.Context, int) (bool, error)
}

type AdminService struct {
	Storage adminStorage
}

var (
	errUsernameUsed     = errors.New("username is already used")
	errIDUsed           = errors.New("id is already used")
	errAdminDoesntExist = errors.New("admin with such username doesn't exist")
	errWrongPassword    = errors.New("wrong password")
)

func (s *AdminService) CreateAdmin(ctx context.Context, admin models.Admin) error {
	ok, err := s.ContainsUsername(ctx, admin.Username)
	if err != nil {
		return err
	}
	if ok {
		return errUsernameUsed
	}

	ok, err = s.ContainsID(ctx, admin.ID)
	if err != nil {
		return err
	}
	if ok {
		return errIDUsed
	}

	err = s.Storage.CreateAdmin(ctx, admin)
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminService) GetAdminByUsername(ctx context.Context, username string) (models.Admin, error) {
	ok, err := s.ContainsUsername(ctx, username)
	if err != nil {
		return models.Admin{}, err
	}
	if !ok {
		return models.Admin{}, errAdminDoesntExist
	}

	return s.Storage.GetAdminByUsername(ctx, username)
}

func (s *AdminService) UpdateAdmin(ctx context.Context, username string, password string, admin models.Admin) error {
	ok, err := s.ContainsUsername(ctx, username)
	if err != nil {
		return err
	}
	if !ok {
		return errAdminDoesntExist
	}

	someAdmin, err := s.Storage.GetAdminByUsername(ctx, username)
	if err != nil {
		return err
	}
	if !someAdmin.CheckPassword(password) {
		return errWrongPassword
	}

	err = s.Storage.UpdateAdmin(ctx, someAdmin.ID, admin)
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminService) DeleteAdmin(ctx context.Context, password string, username string) error {
	ok, err := s.ContainsUsername(ctx, username)
	if err != nil {
		return err
	}
	if !ok {
		return errAdminDoesntExist
	}

	admin, err := s.Storage.GetAdminByUsername(ctx, username)
	if err != nil {
		return err
	}

	if !admin.CheckPassword(password) {
		return errWrongPassword
	}

	err = s.Storage.DeleteAdmin(ctx, username)
	if err != nil {
		return err
	}

	return nil
}

func (s *AdminService) ContainsUsername(ctx context.Context, username string) (bool, error) {
	return s.Storage.ContainsUsername(ctx, username)
}

func (s *AdminService) ContainsID(ctx context.Context, id int) (bool, error) {
	return s.Storage.ContainsID(ctx, id)
}
