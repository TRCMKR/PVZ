package main

import (
	"context"
	"log"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/jsondata"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"
)

const (
	path = "./tests/json_data/data.json"
)

func main() {
	config.InitEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := postgres.NewDB(ctx)
	if err != nil {
		log.Panic(err)
	}

	defer db.GetPool().Close()

	ordersRepo := repository.NewOrderRepo(*db)
	adminsRepo := repository.NewAdminRepo(*db)

	orderStorage, err := jsondata.New(path)
	if err != nil {
		log.Panic("Error: couldn't read json storage", err)
	}

	_ = orderStorage
	_ = ordersRepo
	_ = adminsRepo

	app := web.NewApp(ctx, ordersRepo, adminsRepo)
	if err = app.Run(); err != nil {
		log.Panic("Error: couldn't run app", err)
	}
}
