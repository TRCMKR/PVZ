//go:generate mockgen -source=order_handler.go -destination=../mocks/service/mock_order_service.go -package=service
package web

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	myquery "gitlab.ozon.dev/alexplay1224/homework/internal/query"

	"github.com/Rhymond/go-money"
	"github.com/gorilla/mux"
)

type OrderHandler struct {
	OrderService orderService
}

type orderService interface {
	AcceptOrder(context.Context, int, int, float64, money.Money, time.Time, []models.Packaging) error
	AcceptOrders(context.Context, map[string]models.Order) (int, error)
	ReturnOrder(context.Context, int) error
	ProcessOrders(context.Context, int, []int, string) (int, error)
	UserOrders(context.Context, int, int) ([]models.Order, error)
	Returns(context.Context) ([]models.Order, error)
	GetOrders(context.Context, []myquery.Cond, int, int) ([]models.Order, error)
}

var (
	errNoSuchPackaging   = errors.New("no such packaging")
	errInvalidOrderID    = errors.New("invalid order id")
	errWrongNumberFormat = errors.New("wrong number format")
	errWrongDateFormat   = errors.New("wrong date format")
	errWrongStatusFormat = errors.New("wrong status format")
	errFieldsMissing     = errors.New("missing fields")
	errInvalidUsername   = errors.New("invalid username")
	errWrongJsonFormat   = errors.New("wrong json format")
)

type inputType uint

const (
	numberType inputType = iota
	wordType
	dateType
)

const (
	inputDateAndTimeLayout = "2006.01.02-15:04:05"
	inputDateLayout        = "2006.01.02"
)

func getPackaging(packagingStr string) (models.Packaging, error) {
	var packaging models.Packaging
	if packagingStr != "" {
		packaging = models.GetPackaging(packagingStr)
		if packaging == nil {
			return nil, errNoSuchPackaging
		}
	}

	return packaging, nil
}

type createOrderRequest struct {
	ID             int         `json:"id"`
	UserID         int         `json:"user_id"`
	Weight         float64     `json:"weight"`
	Price          money.Money `json:"price"`
	Packaging      int         `json:"packaging"`
	ExtraPackaging int         `json:"extra_packaging"`
	Status         int         `json:"status"`
	ExpiryDate     time.Time   `json:"expiry_date"`
}

// @Summary		Create an Order
// @Description	Creates a new order.
// @Tags			orders
// @Accept			json
// @Produce		json
// @Param			order	body		models.Order	true	"Order"
// @Success		200		{object}	models.Order
// @Failure		400
// @Router			/orders [post]
func (h *OrderHandler) CreateOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var order = createOrderRequest{}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, errWrongJsonFormat.Error(), http.StatusBadRequest)

		return
	}

	if order.ID == 0 || order.UserID == 0 || order.Weight == 0 || order.Price.IsZero() || order.ExpiryDate.IsZero() {
		http.Error(w, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	packagings := make([]models.Packaging, 0, 2)
	var tmp models.Packaging
	tmp, err = getPackaging(models.GetPackagingName(models.PackagingType(order.Packaging)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	packagings = append(packagings, tmp)
	tmp, err = getPackaging(models.GetPackagingName(models.PackagingType(order.ExtraPackaging)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	packagings = append(packagings, tmp)

	err = h.OrderService.AcceptOrder(ctx, order.ID, order.UserID, order.Weight, order.Price, order.ExpiryDate, packagings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) DeleteOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.Atoi(mux.Vars(r)[orderIDParam])
	if err != nil {
		http.Error(w, errInvalidOrderID.Error(), http.StatusBadRequest)

		return
	}

	err = h.OrderService.ReturnOrder(ctx, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *OrderHandler) parseInt(param string) (int, error) {
	res := 0
	var err error
	if param != "" {
		res, err = strconv.Atoi(param)
		if err != nil {
			return 0, errWrongNumberFormat
		}
	}

	return res, nil
}

func (h *OrderHandler) validateNumberParam(param string) (string, error) {
	if param != "" {
		_, err := strconv.Atoi(param)
		if err != nil {
			return "", errWrongNumberFormat
		}
	}

	return param, nil
}

func (h *OrderHandler) validateWordParam(status string) (string, error) {
	if status == "" {
		return "", nil
	}

	re := regexp.MustCompile("^[a-z]+$")

	if re.MatchString(status) {
		return status, nil
	}

	return "", errWrongStatusFormat
}

func (h *OrderHandler) validateDateParam(param string) (string, error) {
	if param == "" {
		return "", nil
	}

	date, err := time.Parse(inputDateLayout, param)
	if err != nil {
		date, err = time.Parse(inputDateAndTimeLayout, param)
		if err != nil {
			return "", errWrongDateFormat
		}
	}

	return date.Format(time.RFC3339), nil
}

func (h *OrderHandler) getFilters() map[string]inputType {
	return map[string]inputType{
		orderIDParam:         numberType,
		userIDParam:          numberType,
		weightParam:          numberType,
		weightFromParam:      numberType,
		weightToParam:        numberType,
		priceParam:           numberType,
		priceFromParam:       numberType,
		priceToParam:         numberType,
		statusParam:          numberType,
		expiryDateFromParam:  dateType,
		expiryDateToParam:    dateType,
		arrivalDateFromParam: dateType,
		arrivalDateToParam:   dateType,
	}
}

func (h *OrderHandler) getCondMap() map[string]myquery.CondType {
	return map[string]myquery.CondType{
		orderIDParam:         myquery.Equals,
		userIDParam:          myquery.Equals,
		weightParam:          myquery.Equals,
		weightFromParam:      myquery.GreaterEqualThan,
		weightToParam:        myquery.LessEqualThan,
		priceParam:           myquery.Equals,
		priceFromParam:       myquery.GreaterEqualThan,
		priceToParam:         myquery.LessEqualThan,
		statusParam:          myquery.Equals,
		expiryDateFromParam:  myquery.GreaterEqualThan,
		expiryDateToParam:    myquery.LessEqualThan,
		arrivalDateFromParam: myquery.GreaterEqualThan,
		arrivalDateToParam:   myquery.LessEqualThan,
	}
}

func (h *OrderHandler) getColumnMap() map[string]string {
	return map[string]string{
		orderIDParam:         orderIDParam,
		userIDParam:          userIDParam,
		weightParam:          weightParam,
		weightFromParam:      weightParam,
		weightToParam:        weightParam,
		priceParam:           priceParam,
		priceFromParam:       priceParam,
		priceToParam:         priceParam,
		statusParam:          statusParam,
		expiryDateFromParam:  expiryDateParam,
		expiryDateToParam:    expiryDateParam,
		arrivalDateFromParam: arrivalDateParam,
		arrivalDateToParam:   arrivalDateParam,
	}
}

func (h *OrderHandler) getFilterParams(r *http.Request) ([]myquery.Cond, int, int, error) {
	query := r.URL.Query()

	filters := h.getFilters()
	condMap := h.getCondMap()
	columnMap := h.getColumnMap()

	params := make(map[string]string, len(filters))
	var err error
	for k, v := range filters {
		params[k], err = h.validateParam(query.Get(k), v)
		if err != nil {
			return nil, 0, 0, err
		}
	}

	count, err := h.parseInt(query.Get(countParam))
	if err != nil {
		return nil, 0, 0, errWrongNumberFormat
	}
	page, err := h.parseInt(query.Get(pageParam))
	if err != nil {
		return nil, 0, 0, errWrongNumberFormat
	}

	conds := make([]myquery.Cond, 0, len(params))
	for k, v := range params {
		if v == "" {
			continue
		}
		conds = append(conds, myquery.Cond{
			Operator: condMap[k],
			Field:    columnMap[k],
			Value:    v,
		})
	}

	return conds, count, page, nil
}

func (h *OrderHandler) validateParam(value string, inputType inputType) (string, error) {
	switch inputType {
	case numberType:
		return h.validateNumberParam(value)
	case wordType:
		return h.validateWordParam(value)
	case dateType:
		return h.validateDateParam(value)
	default:
		return "", fmt.Errorf("unknown input type: %v", inputType)
	}
}

func (h *OrderHandler) GetOrders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	conds, count, page, err := h.getFilterParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	orders, err := h.OrderService.GetOrders(ctx, conds, count, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	response := struct {
		Count  int            `json:"count"`
		Orders []models.Order `json:"orders"`
	}{
		Count:  len(orders),
		Orders: orders,
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

type processOrderRequest struct {
	UserID   int    `json:"user_id"`
	OrderIDs []int  `json:"order_ids"`
	Action   string `json:"action"`
}

func (h *OrderHandler) UpdateOrders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var processRequest processOrderRequest

	err := json.NewDecoder(r.Body).Decode(&processRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if len(processRequest.OrderIDs) == 0 || processRequest.Action == "" {
		http.Error(w, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	var response = struct {
		Failed int `json:"failed"`
	}{}
	response.Failed, err = h.OrderService.ProcessOrders(ctx, processRequest.UserID,
		processRequest.OrderIDs, processRequest.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
