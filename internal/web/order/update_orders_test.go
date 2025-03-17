//go:build unit

package order

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"gitlab.ozon.dev/alexplay1224/homework/internal/mocks/service"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_UpdateOrders(t *testing.T) {
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
				mockOrderService.EXPECT().ProcessOrders(gomock.Any(), 1, []int{123, 456}, "return").
					Return(0, nil).Times(1)
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
				mockOrderService.EXPECT().ProcessOrders(gomock.Any(), 1, []int{123, 456}, "buy").
					Return(2, errors.New("undefined action")).Times(1)
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

			handler := NewHandler(mockOrderService)

			handler.UpdateOrders(t.Context(), res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)

			if tt.expectedBody != "" {
				assert.JSONEq(t, tt.expectedBody, res.Body.String())
			}
		})
	}
}
