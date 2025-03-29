package facade

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/cache/lru"
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

type AdminFacade struct {
	cache        *lru.Cache[string, models.Admin]
	adminStorage adminStorage
}

func NewAdminFacade(adminStorage adminStorage, capacity int) *AdminFacade {
	return &AdminFacade{
		adminStorage: adminStorage,
		cache:        lru.NewCache[string, models.Admin](capacity),
	}
}

func (f *AdminFacade) CreateAdmin(ctx context.Context, admin models.Admin) error {
	err := f.adminStorage.CreateAdmin(ctx, admin)
	if err != nil {
		return err
	}

	f.cache.Put(admin.Username, admin)

	return nil
}

func (f *AdminFacade) GetAdminByUsername(ctx context.Context, username string) (models.Admin, error) {
	admin, err := f.adminStorage.GetAdminByUsername(ctx, username)
	if err != nil {
		return models.Admin{}, err
	}

	f.cache.Put(username, admin)

	return admin, nil
}

func (f *AdminFacade) UpdateAdmin(ctx context.Context, id int, admin models.Admin) error {
	err := f.adminStorage.UpdateAdmin(ctx, id, admin)
	if err != nil {
		return err
	}

	f.cache.Put(admin.Username, admin)

	return nil
}

func (f *AdminFacade) DeleteAdmin(ctx context.Context, username string) error {
	err := f.adminStorage.DeleteAdmin(ctx, username)
	if err != nil {
		return err
	}

	f.cache.Remove(username)

	return nil
}

func (f *AdminFacade) ContainsUsername(ctx context.Context, username string) (bool, error) {
	if _, ok := f.cache.Get(username); ok {
		return true, nil
	}

	ok, err := f.adminStorage.ContainsUsername(ctx, username)
	if err != nil || !ok {
		return false, err
	}

	admin, _ := f.GetAdminByUsername(ctx, username)

	f.cache.Put(username, admin)

	return ok, nil
}

func (f *AdminFacade) ContainsID(ctx context.Context, id int) (bool, error) {
	return f.adminStorage.ContainsID(ctx, id)
}
