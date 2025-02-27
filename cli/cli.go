package cli

import (
	"bufio"
	"errors"
	"fmt"
	"homework/packaging"
	"os"
	"strings"

	"homework/order"
)

var (
	errScanFailed     = errors.New("scan failed")
	errInvalidCommand = errors.New("invalid command")
)

type storage interface {
	AcceptOrder(orderID int, userID int, weight float64, price float64, expiryDate string,
		packagings []packaging.Packaging) error
	AcceptOrders(path string) (int, int, error)
	ReturnOrder(orderID int) error
	ProcessOrders(userID int, orderIDs []int, action string) (int, error)
	UserOrders(userID int, count int) ([]order.Order, error)
	Returns() []order.Order
	OrderHistory() []order.Order
	Save() error
}

type App struct {
	orderStorage  storage
	stringBuilder strings.Builder
}

func NewApp(appStorage storage) *App {
	return &App{
		orderStorage:  appStorage,
		stringBuilder: strings.Builder{},
	}
}

func (app *App) executeCommand(commands map[string]command,
	inputCommand string, args []string) ([]string, mode, error) {
	if cliCommand, ok := commands[inputCommand]; ok {
		return cliCommand(args)
	}

	return nil, raw, errInvalidCommand
}

func (app *App) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	var result []string
	var md mode
	var err error

	commands := map[string]command{
		"help":          app.help,
		"clear":         app.clearScr,
		"acceptOrder":   app.acceptOrder,
		"acceptOrders":  app.acceptOrders,
		"returnOrder":   app.returnOrder,
		"processOrders": app.processOrders,
		"userOrders":    app.userOrders,
		"returns":       app.returns,
		"orderHistory":  app.orderHistory,
		"exit":          app.exit,
	}

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			return errScanFailed
		}

		line := strings.TrimSpace(scanner.Text())
		args := strings.Split(line, " ")
		inputCommand := args[0]
		args = args[1:]

		if inputCommand == "" {
			continue
		}

		result, md, err = app.executeCommand(commands, inputCommand, args)
		if err != nil {
			fmt.Println("Error:", err)

			continue
		}

		app.draw(result, md)
	}
}
