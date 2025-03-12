package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"
)

func main() {
	config.InitEnv()
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

	app := web.NewApp(ordersRepo, adminsRepo)
	if err = app.Run(ctx); err != nil {
		log.Panic("Error: couldn't run app", err)
	}
}
