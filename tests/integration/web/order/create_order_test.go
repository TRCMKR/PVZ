//go:build integration

package order

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Rhymond/go-money"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	orderServicePkg "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	orderHandlerPkg "gitlab.ozon.dev/alexplay1224/homework/internal/web/order"
	"gitlab.ozon.dev/alexplay1224/homework/tests/integration"
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
		name           string
		args           createOrderRequest
		expectedStatus int
	}{
		{
			name: "Invalid JSON",
			args: createOrderRequest{
				ID: 1,
			},
			expectedStatus: http.StatusBadRequest,
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
			expectedStatus: http.StatusBadRequest,
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
			expectedStatus: http.StatusBadRequest,
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
			expectedStatus: http.StatusBadRequest,
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
			expectedStatus: http.StatusOK,
		},
	}

	ctx := context.Background()
	rootDir, err := config.GetRootDir()
	require.NoError(t, err)
	config.InitEnv(rootDir + "/.env.test")
	cfg := *config.NewConfig()

	connStr, pgContainer, err := integration.InitPostgresContainer(t.Context(), cfg)
	require.NoError(t, err)
	db, err := postgres.NewDB(t.Context(), connStr)
	require.NoError(t, err)
	ordersRepo := repository.NewOrderRepo(*db)
	orderService := orderServicePkg.NewService(ordersRepo)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
		defer db.Close()
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reqBody, err := json.Marshal(tt.args)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(reqBody))
			res := httptest.NewRecorder()

			handler := orderHandlerPkg.NewHandler(orderService)

			handler.CreateOrder(ctx, res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)
			if tt.expectedStatus == http.StatusOK {
				require.Equal(t, "success", res.Body.String())
			}
		})
	}
}
