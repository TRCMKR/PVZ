package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/facade"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/tx_manager"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"
)

func main() {
	err := config.InitEnv(".env")
	if err != nil {
		log.Panic(err)
	}

	cfg := config.NewConfig()

	ctx := context.Background()
	db, err := postgres.NewDB(ctx, cfg.String())
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	tx := tx_manager.NewTxManager(db)

	ordersRepo := repository.NewOrdersRepo(db)
	ordersFacade := facade.NewOrderFacade(ctx, ordersRepo, 10000)
	adminsRepo := repository.NewAdminsRepo(db)
	adminsFacade := facade.NewAdminFacade(adminsRepo, 10000)
	logsRepo := repository.NewLogsRepo(db)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app, err := web.NewApp(ctx, cfg, ordersFacade, adminsFacade, logsRepo, tx, cfg.WorkerCount, cfg.BatchSize, cfg.Timeout)
	if err != nil {
		log.Panic(err)
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
