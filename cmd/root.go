package cmd

import (
	"fmt"
	"os"

	"homework/cmd/cli"
	"homework/models"
	"homework/storage/json_data"
)

func Execute() {

	path := "./storage/json_data/data.json" // подавать при запуске

	var orderStorage models.Storage
	orderStorage, err := json_data.New(path)
	if err != nil {
		fmt.Println("Error: couldn't read json storage", err)
		os.Exit(1)
	}

	err = cli.Start(orderStorage)
	if err != nil {
		fmt.Println("Error: couldn't run app", err)
		os.Exit(1)
	}
}
