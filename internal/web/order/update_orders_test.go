package order

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_UpdateOrders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		requestBody    string
		mockSetup      func(mockOrderService *MockorderService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid request",
			requestBody: `{
                "user_id": 1,
                "id": 123,
                "action": "return"
            }`,
			mockSetup: func(mockOrderService *MockorderService) {
				mockOrderService.EXPECT().ProcessOrder(gomock.Any(), 1, 123, "return").
					Return(nil).Times(1)
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `success`,
		},
		{
			name:           "Invalid JSON",
			requestBody:    `{"user_id": 1, "id": 123, "action": "return"`,
			mockSetup:      func(mockOrderService *MockorderService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing fields",
			requestBody:    `{"user_id": 1, "order_ids":0, "action": ""}`,
			mockSetup:      func(mockOrderService *MockorderService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Service error",
			requestBody: `{
                "user_id": 1,
                "id": 123,
                "action": "buy"
            }`,
			mockSetup: func(mockOrderService *MockorderService) {
				mockOrderService.EXPECT().ProcessOrder(gomock.Any(), 1, 123, "buy").
					Return(errors.New("undefined action")).Times(1)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockOrderService := NewMockorderService(ctrl)
			tt.mockSetup(mockOrderService)

			req := httptest.NewRequest(http.MethodPost, "/orders/update", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			res := httptest.NewRecorder()

			handler := NewHandler(mockOrderService)

			handler.UpdateOrder(t.Context(), res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)

			if tt.expectedBody != "" {
				assert.Equal(t, tt.expectedBody, res.Body.String())
			}
		})
	}
}
