//go:build unit

package order

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.ozon.dev/alexplay1224/homework/internal/mocks/service"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func TestHandler_DeleteOrder(t *testing.T) {
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

			handler := NewHandler(mockOrderService)

			handler.DeleteOrder(t.Context(), res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)
		})
	}
}
