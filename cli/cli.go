package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
	
	"homework/order"
)

var (
	ErrScanFailed     = errors.New("scan failed")
	ErrInvalidCommand = errors.New("invalid command")
)

type storage interface {
	AcceptOrder(orderID int, userID int, expiryDate string) error
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

func (app *App) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	var result []string
	var md mode
	var err error

	var commands = map[string]command{
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
			return ErrScanFailed
		}

		line := strings.TrimSpace(scanner.Text())
		args := strings.Split(line, " ")
		inputCommand := args[0]
		args = args[1:]

		if inputCommand == "" {
			continue
		}

		if cliCommand, ok := commands[inputCommand]; ok {
			result, md, err = cliCommand(args)
		} else {
			err = ErrInvalidCommand
		}
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		app.draw(result, md)
	}
}
