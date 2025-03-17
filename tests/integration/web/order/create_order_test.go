//go:build integration

package order

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/Rhymond/go-money"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	orderServicePkg "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	orderHandlerPkg "gitlab.ozon.dev/alexplay1224/homework/internal/web/order"
	"gitlab.ozon.dev/alexplay1224/homework/tests/integration"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
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
				ExpiryDate:     time.Now().AddDate(1, 0, 0),
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
				ExtraPackaging: -1,
				ExpiryDate:     time.Now().AddDate(1, 0, 0),
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
				ID:             90299,
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

	ctx := context.Background()
	config.InitEnv("../../../../.env.test")
	cfg := *config.NewConfig()

	connStr, pgContainer, err := integration.InitPostgresContainer(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			db, err := postgres.NewDB(ctx, connStr)
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
			fmt.Printf("%s", string(reqBody))
			res := httptest.NewRecorder()

			orderService := orderServicePkg.NewService(ordersRepo)

			handler := orderHandlerPkg.NewHandler(orderService)

			if tt.name == "Correct order" {
				fmt.Print(1)
			}

			handler.CreateOrder(ctx, res, req)

			assert.Equal(t, tt.want, res.Code)
		})
	}
}
