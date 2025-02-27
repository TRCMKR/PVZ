package main

import (
	"log"

	"homework/cli"
	"homework/storage/jsondata"
)

func main() {
	path := "./tests/json_data/data.json"

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
