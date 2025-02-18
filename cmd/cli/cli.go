package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"homework/models"
)

var (
	ErrorScanFailed     = errors.New("scan failed")
	ErrorInvalidCommand = errors.New("invalid command")
)

var (
	orderStorage  models.Storage
	stringBuilder strings.Builder
	//history       []string
)

func Start(s models.Storage) error {
	orderStorage = s

	scanner := bufio.NewScanner(os.Stdin)
	var result []string
	var md mode
	var err error

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			return ErrorScanFailed
		}

		line := strings.TrimSpace(scanner.Text())
		//if len(line) != 0 {
		//	history = append(history, line)
		//}
		args := strings.Split(line, " ")
		inputCommand := args[0]
		args = args[1:]

		if inputCommand == "" {
			//fmt.Println(history)
			continue
		}

		if cliCommand, ok := commands[inputCommand]; ok {
			result, md, err = cliCommand(args)
		} else {
			err = ErrorInvalidCommand
		}
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		draw(result, md)
	}
}
