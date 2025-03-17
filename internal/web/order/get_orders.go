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

func (h *Handler) GetOrders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) getFilters() map[string]inputType {
	return map[string]inputType{
		OrderIDParam:         numberType,
		UserIDParam:          numberType,
		WeightParam:          numberType,
		WeightFromParam:      numberType,
		WeightToParam:        numberType,
		PriceParam:           numberType,
		PriceFromParam:       numberType,
		PriceToParam:         numberType,
		StatusParam:          numberType,
		ExpiryDateFromParam:  dateType,
		ExpiryDateToParam:    dateType,
		ArrivalDateFromParam: dateType,
		ArrivalDateToParam:   dateType,
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

func (h *Handler) validateParam(value string, inputType inputType) (string, error) {
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
