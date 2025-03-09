package web

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service"

	"github.com/Rhymond/go-money"
	"github.com/gorilla/mux"
)

type server struct {
	orderService service.OrderService
	adminService service.AdminService
}

type orderResponse struct {
	Count  int            `json:"count"`
	Orders []models.Order `json:"orders"`
}

var (
	errNoSuchPackaging   = errors.New("no such packaging")
	errInvalidOrderID    = errors.New("invalid order id")
	errWrongNumberFormat = errors.New("wrong number format")
	errWrongDateFormat   = errors.New("wrong date format")
	errWrongStatusFormat = errors.New("wrong status format")
	errFieldsMissing     = errors.New("missing fields")
	errInvalidUsername   = errors.New("invalid username")
)

type orderRequest struct {
	ID             int         `json:"id"`
	UserID         int         `json:"user_id"`
	Weight         float64     `json:"weight"`
	Price          money.Money `json:"price"`
	Packaging      string      `json:"packaging"`
	ExtraPackaging string      `json:"extra_packaging"`
	Status         string      `json:"status"`
	ExpiryDate     time.Time   `json:"expiry_date"`
}

type updateOrdersRequest struct {
	UserID   int    `json:"user_id"`
	OrderIDs []int  `json:"order_ids"`
	Action   string `json:"action"`
}

type updateOrdersResponse struct {
	Failed int `json:"failed"`
}

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

func (s *server) getPackaging(packagingStr string) (models.Packaging, error) {
	var packaging models.Packaging
	if packagingStr != "" {
		packaging = models.GetPackaging(packagingStr)
		if packaging == nil {
			return nil, errNoSuchPackaging
		}
	}

	return packaging, nil
}

func (s *server) CreateOrder(writer http.ResponseWriter, request *http.Request) {
	var order orderRequest

	err := json.NewDecoder(request.Body).Decode(&order)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	if order.ID == 0 || order.UserID == 0 || order.Weight == 0 || order.Price.IsZero() || order.ExpiryDate.IsZero() {
		http.Error(writer, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	packagings := make([]models.Packaging, 0, 2)
	var tmp models.Packaging
	tmp, err = s.getPackaging(order.Packaging)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}
	packagings = append(packagings, tmp)
	tmp, err = s.getPackaging(order.ExtraPackaging)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	packagings = append(packagings, tmp)

	err = s.orderService.AcceptOrder(order.ID, order.UserID, order.Weight, order.Price, order.ExpiryDate, packagings)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (s *server) DeleteOrder(writer http.ResponseWriter, request *http.Request) {
	orderID, err := strconv.Atoi(mux.Vars(request)[orderIDParam])
	if err != nil {
		http.Error(writer, errInvalidOrderID.Error(), http.StatusBadRequest)

		return
	}

	err = s.orderService.ReturnOrder(orderID)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
}

func (s *server) parseInt(param string) (int, error) {
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

func (s *server) validateNumberParam(param string) (string, error) {
	if param != "" {
		_, err := strconv.Atoi(param)
		if err != nil {
			return "", errWrongNumberFormat
		}
	}

	return param, nil
}

func (s *server) validateWordParam(status string) (string, error) {
	if status == "" {
		return "", nil
	}

	re := regexp.MustCompile("^[a-z]+$")

	if re.MatchString(status) {
		return status, nil
	}

	return "", errWrongStatusFormat
}

func (s *server) validateDateParam(param string) (string, error) {
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

func (s *server) getFilterParams(request *http.Request) (map[string]string, int, int, error) {
	query := request.URL.Query()

	filters := map[string]inputType{
		userIDParam:          numberType,
		weightParam:          numberType,
		weightFromParam:      numberType,
		weightToParam:        numberType,
		priceParam:           numberType,
		priceFromParam:       numberType,
		priceToParam:         numberType,
		statusParam:          wordType,
		expiryDateFromParam:  dateType,
		expiryDateToParam:    dateType,
		arrivalDateFromParam: dateType,
		arrivalDateToParam:   dateType,
	}

	params := make(map[string]string, len(filters))
	var err error
	for k, v := range filters {
		switch v {
		case numberType:
			params[k], err = s.validateNumberParam(query.Get(k))
		case wordType:
			params[k], err = s.validateWordParam(query.Get(k))
		case dateType:
			params[k], err = s.validateDateParam(query.Get(k))
		}
		if err != nil {
			return nil, 0, 0, err
		}
	}

	count, err := s.parseInt(query.Get(countParam))
	if err != nil {
		return nil, 0, 0, errWrongNumberFormat
	}
	var page int
	page, err = s.parseInt(query.Get(pageParam))
	if err != nil {
		return nil, 0, 0, errWrongNumberFormat
	}

	return params, count, page, nil
}

func (s *server) GetOrders(writer http.ResponseWriter, request *http.Request) {
	params, count, page, err := s.getFilterParams(request)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	orders := s.orderService.GetOrders(params, count, page)

	response := orderResponse{
		Count:  len(orders),
		Orders: orders,
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(data)
}

func (s *server) UpdateOrders(writer http.ResponseWriter, request *http.Request) {
	var processRequest updateOrdersRequest

	err := json.NewDecoder(request.Body).Decode(&processRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}

	if len(processRequest.OrderIDs) == 0 || processRequest.Action == "" {
		http.Error(writer, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	var response updateOrdersResponse
	response.Failed, err = s.orderService.ProcessOrders(processRequest.UserID,
		processRequest.OrderIDs, processRequest.Action)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	_, _ = writer.Write(data)
}

type createAdminRequest struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (s *server) CreateAdmin(writer http.ResponseWriter, request *http.Request) {
	var createRequest createAdminRequest
	err := json.NewDecoder(request.Body).Decode(&createRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
	}
	if createRequest.ID == 0 || createRequest.Username == "" || createRequest.Password == "" {
		http.Error(writer, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	admin := *models.NewAdmin(createRequest.ID, createRequest.Username, createRequest.Password)

	err = s.adminService.CreateAdmin(admin)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
}

type updateAdminRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

func (s *server) UpdateAdmin(writer http.ResponseWriter, request *http.Request) {
	adminUsername, ok := mux.Vars(request)[adminUsernameParam]
	if !ok {
		http.Error(writer, errInvalidUsername.Error(), http.StatusBadRequest)

		return
	}

	var updateRequest updateAdminRequest
	err := json.NewDecoder(request.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}
	if updateRequest.Password == "" || updateRequest.NewPassword == "" {
		http.Error(writer, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	admin := *models.NewAdmin(0, adminUsername, updateRequest.NewPassword)

	err = s.adminService.UpdateAdmin(adminUsername, updateRequest.Password, admin)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
}

type deleteAdminRequest struct {
	Password string `json:"password"`
}

func (s *server) DeleteAdmin(writer http.ResponseWriter, request *http.Request) {
	adminUsername, ok := mux.Vars(request)[adminUsernameParam]
	if !ok {
		http.Error(writer, errInvalidUsername.Error(), http.StatusBadRequest)

		return
	}
	var deleteRequest deleteAdminRequest
	err := json.NewDecoder(request.Body).Decode(&deleteRequest)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)

		return
	}
	if deleteRequest.Password == "" {
		http.Error(writer, errFieldsMissing.Error(), http.StatusBadRequest)
	}

	err = s.adminService.DeleteAdmin(deleteRequest.Password, adminUsername)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)

		return
	}

	writer.WriteHeader(http.StatusOK)
}
