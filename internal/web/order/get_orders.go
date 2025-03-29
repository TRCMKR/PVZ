package order

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	myquery "gitlab.ozon.dev/alexplay1224/homework/internal/query"
)

type getOrdersResponce struct {
	Count  int            `json:"count"`
	Orders []models.Order `json:"orders"`
}

// GetOrders retrieves a list of orders based on filter parameters
// @Security BasicAuth
// @Summary Get orders with filters
// @Description Retrieves a paginated list of orders based on the provided filter parameters (e.g., count, page, etc.)
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order_id query int false "Order ID"
// @Param user_id query int false "User ID"
// @Param weight query float64 false "Weight of the order"
// @Param weight_from query float64 false "Minimum weight of the order"
// @Param weight_to query float64 false "Maximum weight of the order"
// @Param price query float64 false "Price of the order"
// @Param price_from query float64 false "Minimum price of the order"
// @Param price_to query float64 false "Maximum price of the order"
// @Param status query int false "Status of the order"
// @Param expiry_date_from query string false "Start date of the expiry range" format(date) "2025-03-10T00:00:00Z"
// @Param expiry_date_to query string false "End date of the expiry range" format(date) "2025-03-10T00:00:00Z"
// @Param arrival_date_from query string false "Start date of the arrival range" format(date) "2025-03-10T00:00:00Z"
// @Param arrival_date_to query string false "End date of the arrival range" format(date) "2025-03-10T00:00:00Z"
// @Param count query int false "Number of orders per page"
// @Param page query int false "Page number"
// @Success 200 {object} getOrdersResponce "Success:
// @Failure 400 {string} string "Bad request, invalid parameters"
// @Failure 401 {string} string "Unauthorized"
// @Failure 500 {string} string "Internal server error"
// @Router /orders [get]
func (h *Handler) GetOrders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	conds, count, page, err := h.getFilterParams(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	orders, err := h.OrderService.GetOrders(ctx, conds, count, page)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	response := getOrdersResponce{
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

func (h *Handler) parseInt(param string) (int, error) {
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

func (h *Handler) validateNumberParam(param string) (string, error) {
	if param != "" {
		_, err := strconv.Atoi(param)
		if err != nil {
			return "", errWrongNumberFormat
		}
	}

	return param, nil
}

func (h *Handler) validateWordParam(status string) (string, error) {
	if status == "" {
		return "", nil
	}

	re := regexp.MustCompile("^[a-z]+$")

	if re.MatchString(status) {
		return status, nil
	}

	return "", errWrongStatusFormat
}

func (h *Handler) validateDateParam(param string) (string, error) {
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

func GetFilters() map[string]InputType {
	return map[string]InputType{
		OrderIDParam:         NumberType,
		UserIDParam:          NumberType,
		WeightParam:          NumberType,
		WeightFromParam:      NumberType,
		WeightToParam:        NumberType,
		PriceParam:           NumberType,
		PriceFromParam:       NumberType,
		PriceToParam:         NumberType,
		StatusParam:          NumberType,
		ExpiryDateFromParam:  DateType,
		ExpiryDateToParam:    DateType,
		ArrivalDateFromParam: DateType,
		ArrivalDateToParam:   DateType,
	}
}

func (h *Handler) getCondMap() map[string]myquery.CondType {
	return map[string]myquery.CondType{
		OrderIDParam:         myquery.Equals,
		UserIDParam:          myquery.Equals,
		WeightParam:          myquery.Equals,
		WeightFromParam:      myquery.GreaterEqualThan,
		WeightToParam:        myquery.LessEqualThan,
		PriceParam:           myquery.Equals,
		PriceFromParam:       myquery.GreaterEqualThan,
		PriceToParam:         myquery.LessEqualThan,
		StatusParam:          myquery.Equals,
		ExpiryDateFromParam:  myquery.GreaterEqualThan,
		ExpiryDateToParam:    myquery.LessEqualThan,
		ArrivalDateFromParam: myquery.GreaterEqualThan,
		ArrivalDateToParam:   myquery.LessEqualThan,
	}
}

func (h *Handler) getColumnMap() map[string]string {
	return map[string]string{
		OrderIDParam:         OrderIDParam,
		UserIDParam:          UserIDParam,
		WeightParam:          WeightParam,
		WeightFromParam:      WeightParam,
		WeightToParam:        WeightParam,
		PriceParam:           PriceParam,
		PriceFromParam:       PriceParam,
		PriceToParam:         PriceParam,
		StatusParam:          StatusParam,
		ExpiryDateFromParam:  ExpiryDateParam,
		ExpiryDateToParam:    ExpiryDateParam,
		ArrivalDateFromParam: ArrivalDateParam,
		ArrivalDateToParam:   ArrivalDateParam,
	}
}

func (h *Handler) getFilterParams(r *http.Request) ([]myquery.Cond, int, int, error) {
	query := r.URL.Query()

	filters := GetFilters()
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

	count, err := h.parseInt(query.Get(CountParam))
	if err != nil {
		return nil, 0, 0, errWrongNumberFormat
	}
	page, err := h.parseInt(query.Get(PageParam))
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

func (h *Handler) validateParam(value string, inputType InputType) (string, error) {
	switch inputType {
	case NumberType:
		return h.validateNumberParam(value)
	case WordType:
		return h.validateWordParam(value)
	case DateType:
		return h.validateDateParam(value)
	default:
		return "", fmt.Errorf("unknown input type: %v", inputType)
	}
}
