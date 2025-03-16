//go:build unit

package web

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/Rhymond/go-money"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/alexplay1224/homework/internal/mocks/service"
)

func TestOrderHandler_CreateOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		args      createOrderRequest
		mockSetup func(orderService *service.MockorderService)
		want      int
	}{
		{
			name: "Invalid JSON",
			args: createOrderRequest{
				ID: 1,
			},
			mockSetup: func(orderService *service.MockorderService) {},
			want:      http.StatusBadRequest,
		},
		{
			name: "Missing fields",
			args: createOrderRequest{
				ID:             0,
				UserID:         0,
				Weight:         0,
				Price:          *money.New(0, money.RUB),
				Packaging:      0,
				ExtraPackaging: 0,
				Status:         0,
				ExpiryDate:     time.Time{},
			},
			mockSetup: func(orderService *service.MockorderService) {},
			want:      http.StatusBadRequest,
		},
		{
			name: "Packaging error",
			args: createOrderRequest{
				ID:             -1,
				UserID:         -1,
				Weight:         100,
				Price:          *money.New(1000, money.RUB),
				Packaging:      -1,
				ExtraPackaging: 0,
				Status:         0,
				ExpiryDate:     time.Time{},
			},
			mockSetup: func(orderService *service.MockorderService) {},
			want:      http.StatusBadRequest,
		},
		{
			name: "Not enough weight",
			args: createOrderRequest{
				ID:             123,
				UserID:         2312,
				Weight:         1,
				Price:          *money.New(1000, money.RUB),
				Packaging:      2,
				ExtraPackaging: 3,
				Status:         0,
				ExpiryDate:     time.Now().AddDate(1, 0, 0),
			},
			mockSetup: func(orderService *service.MockorderService) {
				orderService.EXPECT().AcceptOrder(gomock.Any(), gomock.Eq(123), gomock.Eq(2312), gomock.Eq(1.0),
					gomock.Eq(*money.New(1000, money.RUB)),
					gomock.Any(), gomock.Any()).Return(errors.New("not enough weight")).Times(1)
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Correct order",
			args: createOrderRequest{
				ID:             123,
				UserID:         2312,
				Weight:         100,
				Price:          *money.New(1000, money.RUB),
				Packaging:      2,
				ExtraPackaging: 3,
				Status:         0,
				ExpiryDate:     time.Now().AddDate(1, 0, 0),
			},
			mockSetup: func(orderService *service.MockorderService) {
				orderService.EXPECT().AcceptOrder(gomock.Any(), gomock.Eq(123), gomock.Eq(2312), gomock.Eq(100.0),
					gomock.Eq(*money.New(1000, money.RUB)),
					gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			want: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderService := service.NewMockorderService(ctrl)
			tt.mockSetup(mockOrderService)

			reqBody, err := json.Marshal(tt.args)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(reqBody))
			res := httptest.NewRecorder()

			handler := OrderHandler{
				OrderService: mockOrderService,
			}

			handler.CreateOrder(t.Context(), res, req)

			assert.Equal(t, tt.want, res.Code)
		})
	}
}

func TestOrderHandler_GetOrders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		queryParams    map[string]string
		mockSetup      func(orderService *service.MockorderService)
		expectedStatus int
		expectedCount  int
		expectedOrders []models.Order
	}{
		{
			name: "Valid request with filters",
			queryParams: map[string]string{
				"user_id": "123",
				"weight":  "10",
				"count":   "10",
				"page":    "0",
			},
			mockSetup: func(orderService *service.MockorderService) {
				orders := []models.Order{
					{ID: 1, UserID: 123, Weight: 10, Price: *money.New(1000, money.RUB)},
				}
				orderService.EXPECT().GetOrders(gomock.Any(), gomock.Any(), 10, 0).Return(orders, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
			expectedOrders: []models.Order{
				{ID: 1, UserID: 123, Weight: 10, Price: *money.New(1000, money.RUB)},
			},
		},
		{
			name: "Invalid count",
			queryParams: map[string]string{
				"user_id": "123",
				"count":   "invalid",
			},
			mockSetup:      func(orderService *service.MockorderService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
			expectedOrders: nil,
		},
		{
			name: "Invalid filter (invalid weight)",
			queryParams: map[string]string{
				"weight": "not-a-number",
			},
			mockSetup:      func(orderService *service.MockorderService) {},
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
			expectedOrders: nil,
		},
		{
			name: "No filters, page and count provided",
			queryParams: map[string]string{
				"count": "5",
				"page":  "2",
			},
			mockSetup: func(orderService *service.MockorderService) {
				orders := []models.Order{
					{ID: 1, UserID: 123, Weight: 10, Price: *money.New(1000, money.RUB)},
					{ID: 2, UserID: 124, Weight: 20, Price: *money.New(2000, money.RUB)},
				}
				orderService.EXPECT().GetOrders(gomock.Any(), gomock.Any(), 5, 2).Return(orders, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
			expectedOrders: []models.Order{
				{ID: 1, UserID: 123, Weight: 10, Price: *money.New(1000, money.RUB)},
				{ID: 2, UserID: 124, Weight: 20, Price: *money.New(2000, money.RUB)},
			},
		},
		{
			name: "Error from order service",
			queryParams: map[string]string{
				"user_id": "123",
				"count":   "5",
				"page":    "1",
			},
			mockSetup: func(orderService *service.MockorderService) {
				orderService.EXPECT().GetOrders(gomock.Any(), gomock.Any(), 5, 1).Return(nil, errors.New("internal error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
			expectedOrders: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderService := service.NewMockorderService(ctrl)
			tt.mockSetup(mockOrderService)

			req := httptest.NewRequest(http.MethodGet, "/orders", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			res := httptest.NewRecorder()

			handler := &OrderHandler{
				OrderService: mockOrderService,
			}

			handler.GetOrders(t.Context(), res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)

			if tt.expectedStatus == http.StatusOK {
				var response struct {
					Count  int            `json:"count"`
					Orders []models.Order `json:"orders"`
				}
				err := json.NewDecoder(res.Body).Decode(&response)
				if err != nil {
					t.Fatalf("Error decoding response: %v", err)
				}
				assert.Equal(t, tt.expectedCount, response.Count)
				assert.Equal(t, tt.expectedOrders, response.Orders)
			}
		})
	}
}

func TestOrderHandler_DeleteOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		orderIDParam   string
		mockSetup      func(orderService *service.MockorderService)
		expectedStatus int
	}{
		{
			name:         "Valid order ID",
			orderIDParam: "123",
			mockSetup: func(orderService *service.MockorderService) {
				orderService.EXPECT().ReturnOrder(gomock.Any(), 123).Return(nil).Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid order ID format",
			orderIDParam:   "invalid",
			mockSetup:      func(orderService *service.MockorderService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "Error in OrderService.ReturnOrder",
			orderIDParam: "123",
			mockSetup: func(orderService *service.MockorderService) {
				orderService.EXPECT().ReturnOrder(gomock.Any(), 123).Return(errors.New("internal error")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderService := service.NewMockorderService(ctrl)
			tt.mockSetup(mockOrderService)

			req := httptest.NewRequest(http.MethodDelete, "/orders/"+tt.orderIDParam, nil)
			res := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{
				OrderIDParam: tt.orderIDParam,
			})

			handler := &OrderHandler{
				OrderService: mockOrderService,
			}

			handler.DeleteOrder(t.Context(), res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)
		})
	}
}

func TestOrderHandler_UpdateOrders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(mockOrderService *service.MockorderService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid request",
			requestBody: `{
                "user_id": 1,
                "order_ids": [123, 456],
                "action": "return"
            }`,
			mockSetup: func(mockOrderService *service.MockorderService) {
				mockOrderService.EXPECT().ProcessOrders(gomock.Any(), 1, []int{123, 456}, "return").Return(0, nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"failed":0}`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"user_id": 1, "order_ids": [123, 456], "action": "return"`,
			mockSetup:      func(mockOrderService *service.MockorderService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing fields",
			requestBody:    `{"user_id": 1, "order_ids": [], "action": ""}`,
			mockSetup:      func(mockOrderService *service.MockorderService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Service error",
			requestBody: `{
                "user_id": 1,
                "order_ids": [123, 456],
                "action": "buy"
            }`,
			mockSetup: func(mockOrderService *service.MockorderService) {
				mockOrderService.EXPECT().ProcessOrders(gomock.Any(), 1, []int{123, 456}, "buy").Return(2, errors.New("undefined action")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderService := service.NewMockorderService(ctrl)
			tt.mockSetup(mockOrderService)

			req := httptest.NewRequest(http.MethodPost, "/orders/update", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			res := httptest.NewRecorder()

			handler := &OrderHandler{
				OrderService: mockOrderService,
			}

			handler.UpdateOrders(t.Context(), res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, res.Body.String())
			}
		})
	}
}
