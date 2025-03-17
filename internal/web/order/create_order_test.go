//go:build unit

package order

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/mocks/service"

	"github.com/Rhymond/go-money"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandler_CreateOrder(t *testing.T) {
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

			handler := NewHandler(mockOrderService)

			handler.CreateOrder(t.Context(), res, req)

			assert.Equal(t, tt.want, res.Code)
		})
	}
}
