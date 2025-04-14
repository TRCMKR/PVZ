package facade

import (
	"context"

	"github.com/opentracing/opentracing-go"

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

// AdminFacade is a structure for admin facade
type AdminFacade struct {
	cache        *lru.Cache[string, models.Admin]
	adminStorage adminStorage
}

// NewAdminFacade creates an instance for admin facade
func NewAdminFacade(adminStorage adminStorage, capacity int) *AdminFacade {
	return &AdminFacade{
		adminStorage: adminStorage,
		cache:        lru.NewCache[string, models.Admin](capacity),
	}
}

// CreateAdmin creates admin
func (f *AdminFacade) CreateAdmin(ctx context.Context, admin models.Admin) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "adminFacade.CreateAdmin")
	defer span.Finish()

	err := f.adminStorage.CreateAdmin(ctx, admin)
	if err != nil {
		span.SetTag("error", err)

		return err
	}

	f.cache.Put(admin.Username, admin)

	return nil
}

// GetAdminByUsername gets admin by username
func (f *AdminFacade) GetAdminByUsername(ctx context.Context, username string) (models.Admin, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "adminFacade.GetAdminByUsername")
	defer span.Finish()

	if admin, ok := f.cache.Get(username); ok {
		span.SetTag("cache", true)

		return admin, nil
	}

	admin, err := f.adminStorage.GetAdminByUsername(ctx, username)
	if err != nil {
		return models.Admin{}, err
	}

	f.cache.Put(username, admin)

	return admin, nil
}

// UpdateAdmin updates admin by id
func (f *AdminFacade) UpdateAdmin(ctx context.Context, id int, admin models.Admin) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "adminFacade.UpdateAdmin")
	defer span.Finish()

	err := f.adminStorage.UpdateAdmin(ctx, id, admin)
	if err != nil {
		span.SetTag("error", err)

		return err
	}

	f.cache.Put(admin.Username, admin)

	return nil
}

// DeleteAdmin deletes admin
func (f *AdminFacade) DeleteAdmin(ctx context.Context, username string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "adminFacade.DeleteAdmin")
	defer span.Finish()

	err := f.adminStorage.DeleteAdmin(ctx, username)
	if err != nil {
		span.SetTag("error", err)

		return err
	}

	f.cache.Remove(username)

	return nil
}

// ContainsUsername checks if admin by username is present
func (f *AdminFacade) ContainsUsername(ctx context.Context, username string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "adminFacade.ContainsUsername")
	defer span.Finish()

	if _, ok := f.cache.Get(username); ok {
		span.SetTag("cache", true)

		return true, nil
	}

	ok, err := f.adminStorage.ContainsUsername(ctx, username)
	if err != nil || !ok {
		span.SetTag("error", err)

		return false, err
	}

	admin, _ := f.GetAdminByUsername(ctx, username)

	f.cache.Put(username, admin)

	return ok, nil
}

// ContainsID checks if admin by id is present
func (f *AdminFacade) ContainsID(ctx context.Context, id int) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "adminFacade.ContainsID")
	defer span.Finish()

	return f.adminStorage.ContainsID(ctx, id)
}
