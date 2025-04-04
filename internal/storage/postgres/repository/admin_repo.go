package repository

import (
	"context"
	"errors"
	"log"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// AdminsRepo ...
type AdminsRepo struct {
	db database
}

// NewAdminsRepo ...
func NewAdminsRepo(db database) *AdminsRepo {
	return &AdminsRepo{
		db: db,
	}
}

var (
	errCreateAdminFailed        = errors.New("failed to create admin")
	errUpdateAdminFailed        = errors.New("failed to update admin")
	errDeleteAdminFailed        = errors.New("failed to delete admin")
	errGetAdminByUsernameFailed = errors.New("failed to get admin by username")
	errFindingAdmin             = errors.New("could not find admin")
)

// CreateAdmin ...
func (r *AdminsRepo) CreateAdmin(ctx context.Context, admin models.Admin) error {
	_, err := r.db.Exec(ctx, `
							INSERT INTO admins(id, username, password, created_at)
							VALUES ($1, $2, $3, $4)
							`, admin.ID, admin.Username, admin.Password, admin.CreatedAt)
	if err != nil {
		log.Printf("Failed to update order %v", errCreateAdminFailed)

		return errCreateAdminFailed
	}

	return nil
}

// GetAdminByUsername ...
func (r *AdminsRepo) GetAdminByUsername(ctx context.Context, username string) (models.Admin, error) {
	var admin models.Admin
	err := r.db.Get(ctx, &admin, `
								SELECT *
								FROM admins
								WHERE username = $1
								`, username)
	if err != nil {
		log.Printf("Failed to get admin by username %v", errGetAdminByUsernameFailed)

		return models.Admin{}, errGetAdminByUsernameFailed
	}

	return admin, nil
}

// UpdateAdmin ...
func (r *AdminsRepo) UpdateAdmin(ctx context.Context, id int, admin models.Admin) error {
	_, err := r.db.Exec(ctx, `
							UPDATE admins
							SET username = $1, password = $2
							WHERE id = $3
							`, admin.Username, admin.Password, id)
	if err != nil {
		log.Printf("Failed to update admin %v", errUpdateAdminFailed)

		return errUpdateAdminFailed
	}

	return nil
}

// DeleteAdmin ...
func (r *AdminsRepo) DeleteAdmin(ctx context.Context, username string) error {
	_, err := r.db.Exec(ctx, `
							DELETE FROM admins
							WHERE username = $1
							`, username)
	if err != nil {
		log.Printf("Failed to delete admin %v", errDeleteAdminFailed)

		return errDeleteAdminFailed
	}

	return nil
}

// ContainsUsername ...
func (r *AdminsRepo) ContainsUsername(ctx context.Context, username string) (bool, error) {
	var exists bool
	err := r.db.Get(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM admins WHERE username = $1)", username)
	if err != nil {
		log.Printf("Failed to get admin by username %v", errGetAdminByUsernameFailed)

		return false, errGetAdminByUsernameFailed
	}

	return exists, nil
}

// ContainsID ...
func (r *AdminsRepo) ContainsID(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.Get(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM admins WHERE id = $1)", id)
	if err != nil {
		log.Printf("Failed to get admin by id %v", errFindingAdmin)

		return false, errFindingAdmin
	}

	return exists, nil
}
