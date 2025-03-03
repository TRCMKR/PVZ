package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/Rhymond/go-money"
	"github.com/bytedance/sonic"
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
	acceptOrderMinArgCount   = 5
	acceptOrderMaxArgCount   = 7
	acceptOrdersArgCount     = 1
	returnOrderArgCount      = 1
	processOrdersMinArgCount = 3
	userOrdersMinArgCount    = 1
	userOrdersMaxArgCount    = 2
	returnsArgCount          = 0
	orderHistoryArgCount     = 0
)

const (
	inputDateAndTimeLayout = "2006.01.02-15:04:05"
	inputDateLayout        = "2006.01.02"
)

var (
	errTooManyArgs         = errors.New("too many arguments")
	errNotEnoughArgs       = errors.New("not enough arguments")
	errDataNotSaved        = errors.New("data wasn't saved")
	errWrongArgument       = errors.New("wrong argument")
	errWrongPackagingName  = errors.New("wrong packaging name")
	errFileNotOpened       = errors.New("file not opened")
	errDataNotUnmarshalled = errors.New("data not unmarshalled")
	errWrongDateFormat     = errors.New("wrong date format")
)

func checkArgCount(args []string, count int) error {
	if len(args) > count {
		return errTooManyArgs
	} else if len(args) < count {
		return errNotEnoughArgs
	}

	return nil
}

func checkMinMaxArgCount(args []string, minargs int, maxargs int) error {
	if len(args) < minargs {
		return errNotEnoughArgs
	}
	if len(args) > maxargs {
		return errTooManyArgs
	}

	return nil
}

func (a *App) clearScr(args []string) ([]string, mode, error) {
	err := checkArgCount(args, serviceFuncArgCount)
	if err != nil {
		return nil, raw, err
	}

	return []string{clearScreen}, raw, nil
}

func (a *App) help(args []string) ([]string, mode, error) {
	err := checkArgCount(args, serviceFuncArgCount)
	if err != nil {
		return nil, raw, err
	}

	return []string{helpText}, raw, nil
}

func (a *App) exit(args []string) ([]string, mode, error) {
	err := checkArgCount(args, serviceFuncArgCount)
	if err != nil {
		return nil, raw, err
	}

	err = a.orderService.Save()
	if err != nil {
		return nil, raw, errors.Join(errDataNotSaved, err)
	}

	fmt.Println("Exiting...")
	os.Exit(0)

	return nil, raw, nil
}

func stringsToFloats(strings []string) ([]float64, error) {
	floats := make([]float64, len(strings))
	for i, s := range strings {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return nil, err
		}
		floats[i] = f
	}

	return floats, nil
}

func parseInputDate(expiryDate string) (time.Time, error) {
	date, err := time.Parse(inputDateAndTimeLayout, expiryDate)
	if err != nil {
		date, err = time.Parse(inputDateLayout, expiryDate)
		if err != nil {
			return time.Time{}, errWrongDateFormat
		}
	}

	return date, nil
}

func (a *App) acceptOrder(args []string) ([]string, mode, error) {
	err := checkMinMaxArgCount(args, acceptOrderMinArgCount, acceptOrderMaxArgCount)
	if err != nil {
		return nil, raw, err
	}

	intArgs, err := stringsToInts(args[:2])
	if err != nil {
		return nil, raw, errWrongArgument
	}
	orderID, userID := intArgs[0], intArgs[1]

	var floatArgs []float64
	floatArgs, err = stringsToFloats(args[2:4])
	if err != nil {
		return nil, raw, errWrongArgument
	}
	weight, priceFloat := floatArgs[0], floatArgs[1]

	price := *money.NewFromFloat(priceFloat, money.RUB)

	date, err := parseInputDate(args[4])
	if err != nil {
		return nil, raw, err
	}

	args = args[5:]
	packagings := make([]models.Packaging, 0, len(args))
	for _, p := range args {
		tmpPackaging := models.GetPackaging(p)
		if tmpPackaging == nil {
			return nil, raw, errWrongPackagingName
		}
		packagings = append(packagings, tmpPackaging)
	}

	err = a.orderService.AcceptOrder(orderID, userID, weight, price, date, packagings)
	if err != nil {
		return nil, raw, err
	}

	result := "Success: order accepted!"

	return []string{result}, raw, nil
}

func (a *App) getOrdersFromFile(path string) (map[string]models.Order, error) {
	jsonData, err := os.ReadFile(path)
	if err != nil {
		return nil, errFileNotOpened
	}

	data := make(map[string]models.Order)

	err = sonic.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, errDataNotUnmarshalled
	}

	return data, nil
}

func (a *App) acceptOrders(args []string) ([]string, mode, error) {
	err := checkArgCount(args, acceptOrdersArgCount)
	if err != nil {
		return nil, raw, err
	}

	var orders map[string]models.Order
	orders, err = a.getOrdersFromFile(args[0])
	if err != nil {
		return nil, raw, err
	}

	ordersFailed := a.orderService.AcceptOrders(orders)
	orderCount := len(orders)
	var result string
	if ordersFailed > 0 {
		a.stringBuilder.WriteString(fmt.Sprintf("Orders failed: %d/%d", ordersFailed, orderCount))
		result = a.stringBuilder.String()
		a.stringBuilder.Reset()
	} else {
		result = "Success: all orders accepted!"
	}

	return []string{result}, raw, nil
}

func (a *App) returnOrder(args []string) ([]string, mode, error) {
	err := checkArgCount(args, returnOrderArgCount)
	if err != nil {
		return nil, raw, err
	}
	orderID, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, raw, errWrongArgument
	}
	err = a.orderService.ReturnOrder(orderID)
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

func (a *App) formMessage(args []string, ordersFailed int) string {
	if ordersFailed > 0 {
		a.stringBuilder.WriteString("Orders failed: ")
		a.stringBuilder.WriteString(strconv.Itoa(ordersFailed))
	} else {
		a.stringBuilder.WriteString("All orders successfully ")
		switch args[len(args)-1] {
		case "return":
			a.stringBuilder.WriteString("returned")
		case "give":
			a.stringBuilder.WriteString("given")
		}
	}

	result := a.stringBuilder.String()
	a.stringBuilder.Reset()

	return result
}

func (a *App) processOrders(args []string) ([]string, mode, error) {
	if len(args) < processOrdersMinArgCount {
		return nil, raw, errNotEnoughArgs
	}

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, raw, errWrongArgument
	}
	orderIDs, err := stringsToInts(args[1 : len(args)-1])
	if err != nil {
		return nil, raw, errWrongArgument
	}

	ordersFailed, err := a.orderService.ProcessOrders(userID, orderIDs, args[len(args)-1])
	if err != nil {
		return nil, raw, err
	}

	result := a.formMessage(args, ordersFailed)

	return []string{result}, raw, nil
}

func (a *App) parseOptionalArg(args []string) (int, error) {
	count := 0
	var err error
	if len(args) > 0 {
		count, err = strconv.Atoi(args[0])
		if err != nil || count < 1 {
			return 0, errWrongArgument
		}
	}

	return count, nil
}

func (a *App) userOrders(args []string) ([]string, mode, error) {
	err := checkMinMaxArgCount(args, userOrdersMinArgCount, userOrdersMaxArgCount)
	if err != nil {
		return nil, raw, err
	}

	userID, err := strconv.Atoi(args[0])
	if err != nil {
		return nil, raw, errWrongArgument
	}

	var count int
	count, err = a.parseOptionalArg(args[1:])
	if err != nil {
		return nil, raw, err
	}

	orders := a.orderService.UserOrders(userID, count)

	result := make([]string, 0, len(orders))
	for _, someOrder := range orders {
		result = append(result, someOrder.String())
	}

	return result, scrolled, nil
}

func (a *App) returns(args []string) ([]string, mode, error) {
	err := checkArgCount(args, returnsArgCount)
	if err != nil {
		return nil, raw, err
	}

	orders := a.orderService.Returns()
	result := make([]string, 0, len(orders))
	for _, someOrder := range orders {
		result = append(result, someOrder.String())
	}

	return result, paged, nil
}

func (a *App) orderHistory(args []string) ([]string, mode, error) {
	err := checkArgCount(args, orderHistoryArgCount)
	if err != nil {
		return nil, raw, err
	}

	orders := a.orderService.OrderHistory()
	result := make([]string, 0, len(orders))
	for _, someOrder := range orders {
		result = append(result, someOrder.String())
	}

	return result, raw, nil
}
