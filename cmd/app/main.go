package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"
)

func main() {
	if os.Getenv("APP_ENV") == "test" {
		config.InitEnv(".env.test")
	} else {
		config.InitEnv(".env")
	}
	cfg := config.NewConfig()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := postgres.NewDB(ctx, cfg.String())
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	ordersRepo := repository.NewOrderRepo(*db)
	adminsRepo := repository.NewAdminRepo(*db)
	logsRepo := repository.NewLogsRepo(*db)

	app := web.NewApp(ctx, ordersRepo, adminsRepo, logsRepo, 2, 5, 2*time.Second)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Received shutdown signal")
		cancel()
	}()
	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run(ctx)
	}()
	select {
	case <-ctx.Done():
		log.Println("Context canceled, shutting down...")
	case err = <-errCh:
		if err != nil {
			log.Panic("Error: couldn't run app", err)
		}
	}
}
