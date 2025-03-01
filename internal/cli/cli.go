package cli

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"gitlab.ozon.dev/alexplay1224/homework/internal/order"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service"
)

var (
	errScanFailed     = errors.New("scan failed")
	errInvalidCommand = errors.New("invalid command")
)

type storage interface {
	AddOrder(order.Order)
	RemoveOrder(int)
	UpdateOrder(int, order.Order)
	GetByID(int) order.Order
	GetByUserID(int) []order.Order
	GetReturns() []order.Order
	OrderHistory() []order.Order
	Save() error
	Contains(int) bool
}

type App struct {
	orderService  service.OrderService
	stringBuilder strings.Builder
}

func NewApp(appStorage storage) *App {
	return &App{
		orderService: service.OrderService{
			Storage: appStorage,
		},
		stringBuilder: strings.Builder{},
	}
}

func (a *App) executeCommand(commands map[string]command,
	inputCommand string, args []string) ([]string, mode, error) {
	if cliCommand, ok := commands[inputCommand]; ok {
		return cliCommand(args)
	}

	return nil, raw, errInvalidCommand
}

func (a *App) Run() error {
	scanner := bufio.NewScanner(os.Stdin)
	var result []string
	var md mode
	var err error

	commands := map[string]command{
		"help":          a.help,
		"clear":         a.clearScr,
		"acceptOrder":   a.acceptOrder,
		"acceptOrders":  a.acceptOrders,
		"returnOrder":   a.returnOrder,
		"processOrders": a.processOrders,
		"userOrders":    a.userOrders,
		"returns":       a.returns,
		"orderHistory":  a.orderHistory,
		"exit":          a.exit,
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

		result, md, err = a.executeCommand(commands, inputCommand, args)
		if err != nil {
			fmt.Println("Error:", err)

			continue
		}

		a.draw(result, md)
	}
}
