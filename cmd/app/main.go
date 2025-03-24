package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"
)

const (
	workerCount = 1
	batchSize   = 5
	timeout     = 2 * time.Second
)

func main() {
	if os.Getenv("APP_ENV") == "test" {
		config.InitEnv(".env.test")
	} else {
		config.InitEnv(".env")
	}
	cfg := config.NewConfig()

	ctx := context.Background()
	db, err := postgres.NewDB(ctx, cfg.String())
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	ordersRepo := repository.NewOrderRepo(db)
	adminsRepo := repository.NewAdminRepo(db)
	logsRepo := repository.NewLogsRepo(db)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()
	app := web.NewApp(ctx, ordersRepo, adminsRepo, logsRepo, workerCount, batchSize, timeout)
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx)
	}()
	select {
	case <-ctx.Done():
		log.Println("Context canceled, shutting down...")
	case err = <-errCh:
		if err != nil {
			_, file, line, _ := runtime.Caller(1)

			log.Panicf("Error at %s:%d - %v", file, line, err)
		}
	}
}
