//go:build unit

package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"

	"github.com/Rhymond/go-money"
	"github.com/stretchr/testify/assert"
)

type createOrderRequest struct {
	ID             int         `json:"id"`
	UserID         int         `json:"user_id"`
	Weight         float64     `json:"weight"`
	Price          money.Money `json:"price"`
	Packaging      int         `json:"packaging"`
	ExtraPackaging int         `json:"extra_packaging"`
	ExpiryDate     time.Time   `json:"expiry_date"`
}

func TestOrderHandler_CreateOrder(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		args createOrderRequest
		want int
	}{
		{
			name: "Invalid JSON",
			args: createOrderRequest{
				ID: 1,
			},
			want: http.StatusBadRequest,
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
				ExpiryDate:     time.Time{},
			},
			want: http.StatusBadRequest,
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
				ExpiryDate:     time.Time{},
			},
			want: http.StatusBadRequest,
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
				ExpiryDate:     time.Now().AddDate(1, 0, 0),
			},
			want: http.StatusBadRequest,
		},
		{
			name: "Correct order",
			args: createOrderRequest{
				ID:             999,
				UserID:         2312,
				Weight:         100,
				Price:          *money.New(1000, money.RUB),
				Packaging:      2,
				ExtraPackaging: 3,
				ExpiryDate:     time.Now().AddDate(1, 0, 0),
			},
			want: http.StatusOK,
		},
	}
	if os.Getenv("APP_ENV") == "test" {
		config.InitEnv("../../.env.test")
	} else {
		config.InitEnv("../../.env")
	}
	cfg := config.NewConfig()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			db, err := postgres.NewDB(ctx, cfg.String())
			if err != nil {
				log.Panic(err)
			}

			defer db.Close()

			ordersRepo := repository.NewOrderRepo(*db)

			reqBody, err := json.Marshal(tt.args)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(reqBody))
			res := httptest.NewRecorder()

			orderService := service.OrderService{
				Storage: ordersRepo,
			}

			handler := web.OrderHandler{
				OrderService: &orderService,
			}

			handler.CreateOrder(ctx, res, req)

			if res.Code == http.StatusOK {
				handler.DeleteOrder(ctx, res, req)
			}

			assert.Equal(t, tt.want, res.Code)
		})
	}
}
