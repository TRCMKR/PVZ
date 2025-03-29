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
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/facade"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/tx_manager"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"
)

const (
	workerCount = 2
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

	tx := tx_manager.NewTxManager(db)

	ordersRepo := repository.NewOrderRepo(tx)
	ordersFacade := facade.NewOrderFacade(ctx, ordersRepo, 10000)
	adminsRepo := repository.NewAdminRepo(db)
	adminsFacade := facade.NewAdminFacade(adminsRepo, 10000)
	logsRepo := repository.NewLogsRepo(db)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app, err := web.NewApp(ctx, ordersFacade, adminsFacade, logsRepo, tx, workerCount, batchSize, timeout)
	if err != nil {
		log.Fatal(err)
	}

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
