//go:build unit

package order

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.ozon.dev/alexplay1224/homework/internal/mocks/service"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/Rhymond/go-money"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_GetOrders(t *testing.T) {
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

			handler := NewHandler(mockOrderService)

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
