package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

const (
	helpText = `Usage:
  [command] arguments...
Available Commands:
  help               prints this help message
  clear              clears screen
  acceptOrder        accepts arg1 order for arg2 user with arg3 expiration date
  acceptOrders       accepts orders from other json
  returnOrder        returns to courier expired, not given arg1 order
  processOrders      gives to (or makes returns from) arg1 user arg... orders
  userOrders         shows last arg2 (optional) orders by arg1 user
  returns            shows all returned orders
  orderHistory       shows order history
  exit               ends program`

	clearScreen = "\033[H\033[2J"
)

type command func(args []string) ([]string, mode, error)

type mode int

const (
	raw mode = iota
	paged
	scrolled
)

const (
	serviceFuncArgCount      = 0
	acceptOrderArgCount      = 3
	acceptOrdersArgCount     = 1
	returnOrderArgCount      = 1
	processOrdersMinArgCount = 3
	userOrdersMinArgCount    = 1
	userOrdersMaxArgCount    = 2
	returnsArgCount          = 0
	orderHistoryArgCount     = 0
)

var (
	ErrTooManyArgs   = errors.New("too many arguments")
	ErrNotEnoughArgs = errors.New("not enough arguments")
	ErrDataNotSaved  = errors.New("data wasn't saved")
	ErrWrongArgument = errors.New("wrong argument")
)

func checkArgCount(args []string, count int) error {
	if len(args) > count {
		return ErrTooManyArgs
	} else if len(args) < count {
		return ErrNotEnoughArgs
	}

	return nil
}

func (app *App) clearScr(args []string) ([]string, mode, error) {
	err := checkArgCount(args, serviceFuncArgCount)
	if err != nil {
		return nil, raw, err
	}

	return []string{clearScreen}, raw, nil
}

func (app *App) help(args []string) ([]string, mode, error) {
	err := checkArgCount(args, serviceFuncArgCount)
	if err != nil {
		return nil, raw, err
	}

	return []string{helpText}, raw, nil
}

func (app *App) exit(args []string) ([]string, mode, error) {
	err := checkArgCount(args, serviceFuncArgCount)
	if err != nil {
		return nil, raw, err
	}

	err = app.orderStorage.Save()
	if err != nil {
		return nil, raw, errors.Join(ErrDataNotSaved, err)
	}

	fmt.Println("Exiting...")
	os.Exit(0)
	return nil, raw, nil
}

func (app *App) acceptOrder(args []string) ([]string, mode, error) {
	err := checkArgCount(args, acceptOrderArgCount)
	if err != nil {
		return nil, raw, err
	}

	orderID, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, raw, ErrWrongArgument
	}
	userID, err := strconv.Atoi(args[1])
	if err != nil {
		return nil, raw, ErrWrongArgument
	}

	err = app.orderStorage.AcceptOrder(orderID, userID, args[2])
	if err != nil {
		return nil, raw, err
	}

	result := "Success: order accepted!"
	return []string{result}, raw, nil
}

func (app *App) acceptOrders(args []string) ([]string, mode, error) {
	err := checkArgCount(args, acceptOrdersArgCount)
	if err != nil {
		return nil, raw, err
	}

	var ordersFailed int
	var orderCount int
	ordersFailed, orderCount, err = app.orderStorage.AcceptOrders(args[0])
	if err != nil {
		return nil, raw, err
	}

	var result string
	if ordersFailed > 0 {
		app.stringBuilder.WriteString(fmt.Sprintf("Orders failed: %d/%d", ordersFailed, orderCount))
		result = app.stringBuilder.String()
		app.stringBuilder.Reset()
	} else {
		result = "Success: all orders accepted!"
	}

	return []string{result}, raw, nil
}

func (app *App) returnOrder(args []string) ([]string, mode, error) {
	err := checkArgCount(args, returnOrderArgCount)
	if err != nil {
		return nil, raw, err
	}
	orderID, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, raw, ErrWrongArgument
	}
	err = app.orderStorage.ReturnOrder(orderID)
	if err != nil {
		return nil, raw, err
	}

	result := "Success: order returned!"
	return []string{result}, raw, nil
}

func stringsToInts(strings []string) ([]int, error) {
	result := make([]int, len(strings))
	for i, s := range strings {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		result[i] = n
	}

	return result, nil
}

func (app *App) processOrders(args []string) ([]string, mode, error) {
	if len(args) < processOrdersMinArgCount {
		return nil, raw, ErrNotEnoughArgs
	}

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, raw, ErrWrongArgument
	}
	orderIDs, err := stringsToInts(args[1 : len(args)-1])
	if err != nil {
		return nil, raw, ErrWrongArgument
	}

	ordersFailed, err := app.orderStorage.ProcessOrders(userID, orderIDs, args[len(args)-1])
	if err != nil {
		return nil, raw, err
	}

	if ordersFailed > 0 {
		app.stringBuilder.WriteString("Orders failed: ")
		app.stringBuilder.WriteString(strconv.Itoa(ordersFailed))
	} else {
		app.stringBuilder.WriteString("All orders successfully ")
		switch args[len(args)-1] {
		case "return":
			app.stringBuilder.WriteString("returned")
		case "give":
			app.stringBuilder.WriteString("given")
		}
	}
	result := app.stringBuilder.String()
	app.stringBuilder.Reset()

	return []string{result}, raw, nil
}

func (app *App) userOrders(args []string) ([]string, mode, error) {
	if len(args) < userOrdersMinArgCount {
		return nil, raw, ErrNotEnoughArgs
	}
	if len(args) > userOrdersMaxArgCount {
		return nil, raw, ErrTooManyArgs
	}

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, raw, ErrWrongArgument
	}
	args = args[1:]
	count := 0
	if len(args) > 0 {
		count, err = strconv.Atoi(args[0])
		if err != nil || count < 1 {
			return nil, raw, ErrWrongArgument
		}
	}

	orders, err := app.orderStorage.UserOrders(userID, count)
	if err != nil {
		return nil, raw, err
	}

	result := make([]string, 0, len(orders))
	for _, order := range orders {
		result = append(result, order.String())
	}

	return result, scrolled, nil
}

func (app *App) returns(args []string) ([]string, mode, error) {
	err := checkArgCount(args, returnsArgCount)
	if err != nil {
		return nil, raw, err
	}

	orders := app.orderStorage.Returns()
	var result []string
	for _, order := range orders {
		result = append(result, order.String())
	}

	return result, paged, nil
}

func (app *App) orderHistory(args []string) ([]string, mode, error) {
	err := checkArgCount(args, orderHistoryArgCount)
	if err != nil {
		return nil, raw, err
	}

	var result []string
	orders := app.orderStorage.OrderHistory()
	for _, order := range orders {
		result = append(result, order.String())
	}

	return result, raw, nil
}
