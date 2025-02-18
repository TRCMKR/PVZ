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

var commands = map[string]command{
	"help":          help,
	"clear":         clearScr,
	"acceptOrder":   acceptOrder,
	"acceptOrders":  acceptOrders,
	"returnOrder":   returnOrder,
	"processOrders": processOrders,
	"userOrders":    userOrders,
	"returns":       returns,
	"orderHistory":  orderHistory,
	"exit":          exit,
}

var (
	ErrorTooManyArgs   = errors.New("too many arguments")
	ErrorNotEnoughArgs = errors.New("not enough arguments")
	ErrorDataNotSaved  = errors.New("data wasn't saved")
)

func checkArgsCount(args []string, count int) error {
	if len(args) > count {
		return ErrorTooManyArgs
	} else if len(args) < count {
		return ErrorNotEnoughArgs
	}

	return nil
}

func clearScr(args []string) ([]string, mode, error) {
	err := checkArgsCount(args, 0)
	if err != nil {
		return nil, raw, err
	}

	return []string{clearScreen}, raw, nil
}

func help(args []string) ([]string, mode, error) {
	err := checkArgsCount(args, 0)
	if err != nil {
		return nil, raw, err
	}

	return []string{helpText}, raw, nil
}

func exit(args []string) ([]string, mode, error) {
	err := checkArgsCount(args, 0)
	if err != nil {
		return nil, raw, err
	}

	err = orderStorage.Save()
	if err != nil {
		return nil, raw, errors.Join(ErrorDataNotSaved, err)
	}

	fmt.Println("Exiting...")
	os.Exit(0)
	return nil, raw, nil
}

func acceptOrder(args []string) ([]string, mode, error) {
	err := checkArgsCount(args, 3)
	if err != nil {
		return nil, raw, err
	}

	err = orderStorage.AcceptOrder(args[0], args[1], args[2])
	if err != nil {
		return nil, raw, err
	}

	result := "Success: order accepted!\n"
	return []string{result}, raw, nil
}

func acceptOrders(args []string) ([]string, mode, error) {
	err := checkArgsCount(args, 1)
	if err != nil {
		return nil, raw, err
	}

	var ordersFailed int
	ordersFailed, err = orderStorage.AcceptOrders(args[0])
	if err != nil {
		return nil, raw, err
	}

	var result string
	if ordersFailed > 0 {
		stringBuilder.WriteString("Orders failed: ")
		stringBuilder.WriteString(strconv.Itoa(ordersFailed))
		stringBuilder.WriteByte('\n')
		result = stringBuilder.String()
		stringBuilder.Reset()
	} else {
		result = "Success: all orders accepted!\n"
	}

	return []string{result}, raw, nil
}

func returnOrder(args []string) ([]string, mode, error) {
	err := checkArgsCount(args, 1)
	if err != nil {
		return nil, raw, err
	}
	err = orderStorage.ReturnOrder(args[0])
	if err != nil {
		return nil, raw, err
	}

	result := "Success: order returned!\n"
	return []string{result}, raw, nil
}

func processOrders(args []string) ([]string, mode, error) {
	if len(args) < 3 {
		return nil, raw, ErrorNotEnoughArgs
	}

	ordersFailed, err := orderStorage.ProcessOrders(args[0], args[1:len(args)-1], args[len(args)-1])
	if err != nil {
		return nil, raw, err
	}

	if ordersFailed > 0 {
		stringBuilder.WriteString("Orders failed: ")
		stringBuilder.WriteString(strconv.Itoa(ordersFailed))
	} else {
		stringBuilder.WriteString("All orders successfully ")
		if args[len(args)-1] == "return" {
			stringBuilder.WriteString("returned")
		} else {
			stringBuilder.WriteString("given")
		}
	}
	stringBuilder.WriteByte('\n')
	result := stringBuilder.String()
	stringBuilder.Reset()

	return []string{result}, raw, nil
}

func userOrders(args []string) ([]string, mode, error) {
	if len(args) < 1 {
		return nil, raw, ErrorNotEnoughArgs
	} else if len(args) > 2 {
		return nil, raw, ErrorTooManyArgs
	}

	orders, err := orderStorage.UserOrders(args...)
	if err != nil {
		return nil, raw, err
	}

	var result []string
	for _, order := range orders {
		result = append(result, order.String())
	}

	return result, scrolled, nil
}

func returns(args []string) ([]string, mode, error) {
	err := checkArgsCount(args, 0)
	if err != nil {
		return nil, raw, err
	}

	orders := orderStorage.Returns()
	var result []string
	for _, order := range orders {
		result = append(result, order.String())
	}

	return result, paged, nil
}

func orderHistory(args []string) ([]string, mode, error) {
	err := checkArgsCount(args, 0)
	if err != nil {
		return nil, raw, err
	}

	var result []string
	orders := orderStorage.OrderHistory()
	for _, o := range orders {
		result = append(result, o.String())
	}

	return result, raw, nil
}
