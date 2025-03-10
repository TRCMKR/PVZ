package main

import (
	"context"
	"fmt"
	"log"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"
)

func generateDsn() string {
	host := config.GetDBHost()
	port := config.GetDBPort()
	username := config.GetDBUsername()
	password := config.GetDBPassword()
	dbname := config.GetDBName()

	if host == "" || port == "" || username == "" || password == "" || dbname == "" {
		log.Fatal("Database configuration missing: one or more required fields are empty.")
	}

	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, dbname)
}

func main() {
	config.InitEnv()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db, err := postgres.NewDB(ctx, generateDsn())
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
