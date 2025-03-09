package service

import (
	"context"
	"errors"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

type adminStorage interface {
	CreateAdmin(context.Context, models.Admin)
	GetAdminByUsername(context.Context, string) models.Admin
	UpdateAdmin(context.Context, int, models.Admin)
	DeleteAdmin(context.Context, string)
	ContainsUsername(context.Context, string) bool
	ContainsID(context.Context, int) bool
}

type AdminService struct {
	Storage adminStorage
	Ctx     context.Context
}

var (
	errUsernameUsed     = errors.New("username is already used")
	errIDUsed           = errors.New("id is already used")
	errAdminDoesntExist = errors.New("admin with such username doesn't exist")
	errWrongPassword    = errors.New("wrong password")
)

func (s *AdminService) CreateAdmin(admin models.Admin) error {
	if s.ContainsUsername(admin.Username) {
		return errUsernameUsed
	}
	if s.ContainsID(admin.ID) {
		return errIDUsed
	}

	s.Storage.CreateAdmin(s.Ctx, admin)

	return nil
}

func (s *AdminService) GetAdminByUsername(username string) (models.Admin, error) {
	if !s.ContainsUsername(username) {
		return models.Admin{}, errAdminDoesntExist
	}

	return s.Storage.GetAdminByUsername(s.Ctx, username), nil
}

func (s *AdminService) UpdateAdmin(username string, password string, admin models.Admin) error {
	if !s.ContainsUsername(username) {
		return errAdminDoesntExist
	}

	someAdmin := s.Storage.GetAdminByUsername(s.Ctx, username)
	if !someAdmin.CheckPassword(password) {
		return errWrongPassword
	}

	s.Storage.UpdateAdmin(s.Ctx, someAdmin.ID, admin)

	return nil
}

func (s *AdminService) DeleteAdmin(password string, username string) error {
	if !s.ContainsUsername(username) {
		return errAdminDoesntExist
	}

	admin := s.Storage.GetAdminByUsername(s.Ctx, username)

	if !admin.CheckPassword(password) {
		return errWrongPassword
	}

	s.Storage.DeleteAdmin(s.Ctx, username)

	return nil
}

func (s *AdminService) ContainsUsername(username string) bool {
	return s.Storage.ContainsUsername(s.Ctx, username)
}

func (s *AdminService) ContainsID(id int) bool {
	return s.Storage.ContainsID(s.Ctx, id)
}
