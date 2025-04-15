package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/zap"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/facade"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/tx_manager"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web/grpc"
)

func main() {
	err := config.InitEnv(".env")
	if err != nil {
		log.Panic(err)
	}

	cfg := config.NewConfig()

	_, closer, err := config.InitTracer("grpc-app")
	if err != nil {
		log.Panic("cannot start tracer", err)
	}
	defer closer.Close()

	ctx := context.Background()
	db, err := postgres.NewDB(ctx, cfg.String())
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()

	logger, err := zap.NewProduction()
	if err != nil {
		log.Panic("cannot init zap logger", err)
	}
	defer logger.Sync()

	tx := tx_manager.NewTxManager(db)

	ordersRepo := repository.NewOrdersRepo(logger.With(
		zap.String("layer", "orders repo"),
	), db)
	ordersFacade := facade.NewOrderFacade(ctx, ordersRepo, 10000)

	adminsRepo := repository.NewAdminsRepo(logger.With(
		zap.String("layer", "admins repo"),
	), db)
	adminsFacade := facade.NewAdminFacade(adminsRepo, 10000)

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	app := grpc.NewServer(logger, ordersFacade, adminsFacade, tx)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx, cfg, logger)
	}()

	select {
	case <-ctx.Done():
		log.Println("Context canceled, shutting down...")
	case err = <-errCh:
		if err != nil {
			_, file, line, _ := runtime.Caller(1)

			logger.Panic(fmt.Sprintf("Error at %s:%d - %v", file, line, err))
		}
	}
}
