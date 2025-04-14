package repository

import (
	"context"
	"errors"

	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// AdminsRepo is a structure for admins repo
type AdminsRepo struct {
	db     database
	logger *zap.Logger
}

// NewAdminsRepo creates an instance of admins repo
func NewAdminsRepo(logger *zap.Logger, db database) *AdminsRepo {
	return &AdminsRepo{
		db:     db,
		logger: logger,
	}
}

var (
	errCreateAdminFailed        = errors.New("failed to create admin")
	errUpdateAdminFailed        = errors.New("failed to update admin")
	errDeleteAdminFailed        = errors.New("failed to delete admin")
	errGetAdminByUsernameFailed = errors.New("failed to get admin by username")
	errFindingAdmin             = errors.New("could not find admin")
)

// CreateAdmin creates admin
func (r *AdminsRepo) CreateAdmin(ctx context.Context, admin models.Admin) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.CreateAdmin")
	defer span.Finish()

	_, err := r.db.Exec(ctx, `
							INSERT INTO admins(id, username, password, created_at)
							VALUES ($1, $2, $3, $4)
							`, admin.ID, admin.Username, admin.Password, admin.CreatedAt)
	if err != nil {
		r.logger.Error("failed to insert admin",
			zap.Int("id", admin.ID),
			zap.String("username", admin.Username),
			zap.Error(err),
		)
		span.SetTag("error", errCreateAdminFailed)

		return errCreateAdminFailed
	}

	return nil
}

// GetAdminByUsername gets admin by username
func (r *AdminsRepo) GetAdminByUsername(ctx context.Context, username string) (models.Admin, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.GetAdminByUsername")
	defer span.Finish()

	var admin models.Admin
	err := r.db.Get(ctx, &admin, `
								SELECT *
								FROM admins
								WHERE username = $1
								`, username)
	if err != nil {
		r.logger.Error("failed to insert admin",
			zap.String("username", username),
			zap.Error(err),
		)
		span.SetTag("error", errGetAdminByUsernameFailed)

		return models.Admin{}, errGetAdminByUsernameFailed
	}

	return admin, nil
}

// UpdateAdmin updates admin
func (r *AdminsRepo) UpdateAdmin(ctx context.Context, id int, admin models.Admin) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.UpdateAdmin")
	defer span.Finish()

	_, err := r.db.Exec(ctx, `
							UPDATE admins
							SET username = $1, password = $2
							WHERE id = $3
							`, admin.Username, admin.Password, id)
	if err != nil {
		r.logger.Error("failed to update admin",
			zap.Int("id", id),
			zap.String("username", admin.Username),
			zap.Error(err),
		)
		span.SetTag("error", errUpdateAdminFailed)

		return errUpdateAdminFailed
	}

	return nil
}

// DeleteAdmin deletes admin
func (r *AdminsRepo) DeleteAdmin(ctx context.Context, username string) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.DeleteAdmin")
	defer span.Finish()

	_, err := r.db.Exec(ctx, `
							DELETE FROM admins
							WHERE username = $1
							`, username)
	if err != nil {
		r.logger.Error("failed to delete admin",
			zap.String("username", username),
			zap.Error(err),
		)
		span.SetTag("error", errDeleteAdminFailed)

		return errDeleteAdminFailed
	}

	return nil
}

// ContainsUsername checks if admin by username is present
func (r *AdminsRepo) ContainsUsername(ctx context.Context, username string) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.ContainsUsername")
	defer span.Finish()

	var exists bool
	err := r.db.Get(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM admins WHERE username = $1)", username)
	if err != nil {
		r.logger.Error("failed to check if admin exists",
			zap.String("username", username),
			zap.Error(err),
		)
		span.SetTag("error", errGetAdminByUsernameFailed)

		return false, errGetAdminByUsernameFailed
	}

	return exists, nil
}

// ContainsID checks if admin by id is present
func (r *AdminsRepo) ContainsID(ctx context.Context, id int) (bool, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "repo.ContainsID")
	defer span.Finish()

	var exists bool
	err := r.db.Get(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM admins WHERE id = $1)", id)
	if err != nil {
		r.logger.Error("failed to check if admin exists",
			zap.Int("id", id),
			zap.Error(err),
		)
		span.SetTag("error", errFindingAdmin)

		return false, errFindingAdmin
	}

	return exists, nil
}
