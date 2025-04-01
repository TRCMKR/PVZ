package integration

import (
	"context"
	"database/sql"
	"io"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go/wait"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"

	// lib/pg ...
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// InitPostgresContainer ...
func InitPostgresContainer(ctx context.Context, cfg config.Config) (string, *postgres.PostgresContainer, error) {
	emptyLogger := log.New(io.Discard, "", 0)

	pgContainer, err := postgres.Run(ctx, "postgres:14-alpine",
		postgres.WithDatabase(cfg.DBName()),
		postgres.WithUsername(cfg.Username()),
		postgres.WithPassword(cfg.Password()),
		testcontainers.WithLogger(emptyLogger),
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

	goose.SetLogger(emptyLogger)
	err = goose.SetDialect("postgres")
	if err != nil {
		return "", nil, err
	}
	rootDir, err := config.GetRootDir()
	if err != nil {
		return "", nil, err
	}
	if err = goose.Up(db, rootDir+"/migrations"); err != nil {
		log.Panicf("failed to run migrations: %v", err)
	}

	return connStr, pgContainer, nil
}
