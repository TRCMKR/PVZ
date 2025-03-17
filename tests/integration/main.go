//go:build integration

package integration

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"
	"gitlab.ozon.dev/alexplay1224/homework/internal/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
)

func InitPostgresContainer(ctx context.Context, cfg config.Config) (string, *pgcontainer.PostgresContainer, error) {
	pgContainer, err := pgcontainer.Run(ctx, "postgres:14-alpine",
		pgcontainer.WithDatabase(cfg.DBName()),
		pgcontainer.WithUsername(cfg.Username()),
		pgcontainer.WithPassword(cfg.Password()),
		testcontainers.WithLogger(log.Default()),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return "", nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return "", nil, err
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return "", nil, err
	}
	defer db.Close()

	time.Sleep(2 * time.Second)

	goose.SetDialect("postgres")
	if err := goose.Up(db, "../../../../migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return connStr, pgContainer, nil
}
