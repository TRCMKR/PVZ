package repository

import (
	"context"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
)

type AdminRepo struct {
	db postgres.Database
}

func NewAdminRepo(db postgres.Database) *AdminRepo {
	return &AdminRepo{
		db: db,
	}
}

func (r *AdminRepo) CreateAdmin(ctx context.Context, admin models.Admin) {
	_, _ = r.db.Exec(ctx, `INSERT INTO admins(id, username, password, created_at)
		VALUES ($1, $2, $3, $4)`, admin.ID, admin.Username, admin.Password, admin.CreatedAt)
}

func (r *AdminRepo) GetAdminByUsername(ctx context.Context, username string) models.Admin {
	var admin models.Admin
	_ = r.db.Get(ctx, &admin, "SELECT * FROM admins WHERE username = $1", username)

	return admin
}

func (r *AdminRepo) UpdateAdmin(ctx context.Context, id int, admin models.Admin) {
	_, _ = r.db.Exec(ctx, `UPDATE admins
	SET username = $1, password = $2
	WHERE id = $3`, admin.Username, admin.Password, id)
}

func (r *AdminRepo) DeleteAdmin(ctx context.Context, username string) {
	_, _ = r.db.Exec(ctx, "DELETE FROM admins WHERE username = $1", username)
}

func (r *AdminRepo) ContainsUsername(ctx context.Context, username string) bool {
	var exists bool
	_ = r.db.Get(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM admins WHERE username = $1)", username)

	return exists
}

func (r *AdminRepo) ContainsID(ctx context.Context, id int) bool {
	var exists bool
	_ = r.db.Get(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM admins WHERE id = $1)", id)

	return exists
}
