package order

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestHandler_DeleteOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		orderIDParam   string
		mockSetup      func(orderService *MockorderService)
		expectedStatus int
	}{
		{
			name:         "Valid order ID",
			orderIDParam: "123",
			mockSetup: func(orderService *MockorderService) {
				orderService.EXPECT().ReturnOrder(gomock.Any(), 123).Return(nil).Times(1)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid order ID format",
			orderIDParam:   "invalid",
			mockSetup:      func(orderService *MockorderService) {},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "Error in OrderService.ReturnOrder",
			orderIDParam: "123",
			mockSetup: func(orderService *MockorderService) {
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

			mockOrderService := NewMockorderService(ctrl)
			tt.mockSetup(mockOrderService)

			req := httptest.NewRequest(http.MethodDelete, "/orders/"+tt.orderIDParam, nil)
			res := httptest.NewRecorder()
			req = mux.SetURLVars(req, map[string]string{
				OrderIDParam: tt.orderIDParam,
			})

			handler := NewHandler(mockOrderService)

			handler.DeleteOrder(t.Context(), res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)
		})
	}
}
