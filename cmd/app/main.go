package main

import (
	"log"

	"gitlab.ozon.dev/alexplay1224/homework/internal/cli"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/jsondata"
)

const (
	path = "./tests/json_data/data.json"
)

func main() {
	orderStorage, err := jsondata.New(path)
	if err != nil {
		log.Fatal("Error: couldn't read json storage", err)
	}

	app := cli.NewApp(orderStorage)
	err = app.Run()
	if err != nil {
		log.Fatal("Error: couldn't run app", err)
	}
}
