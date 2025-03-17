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
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
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

func TestOrderHandler_UpdateOrders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		failed         int
	}{
		{
			name: "Update existing order",
			requestBody: `{
                "user_id": 123,
                "order_ids": [1],
                "action": "give"
            }`,
			expectedStatus: http.StatusOK,
			failed:         0,
		},
		{
			name: "Update non-existing order",
			requestBody: `{
                "user_id": 123,
                "order_ids": [1232],
                "action": "give"
            }`,
			expectedStatus: http.StatusOK,
			failed:         1,
		},
	}

	ctx := context.Background()
	config.InitEnv("../../../../.env.test")
	cfg := *config.NewConfig()

	connStr, pgContainer, err := integration.InitPostgresContainer(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}

	db, err := postgres.NewDB(ctx, connStr)
	if err != nil {
		log.Panic(err)
	}

	defer db.Close()

	ordersRepo := repository.NewOrderRepo(*db)

	orders := []models.Order{
		{
			ID:          1,
			UserID:      123,
			Weight:      10,
			Price:       *money.New(1000, money.RUB),
			Status:      1,
			ArrivalDate: time.Now(),
			LastChange:  time.Now(),
			ExpiryDate:  time.Now().AddDate(1, 0, 0),
		},
	}

	for _, order := range orders {
		ordersRepo.AddOrder(ctx, order)
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
			orderService := orderServicePkg.NewService(ordersRepo)

			reqBody := []byte(tt.requestBody)

			req := httptest.NewRequest(http.MethodPost, "/orders/process", bytes.NewReader(reqBody))
			res := httptest.NewRecorder()
			handler := orderHandlerPkg.NewHandler(orderService)

			handler.UpdateOrders(ctx, res, req)

			var response = struct {
				Failed int `json:"failed"`
			}{}

			assert.Equal(t, tt.expectedStatus, res.Code)
			if tt.expectedStatus == http.StatusOK {
				require.NoError(t, json.Unmarshal(res.Body.Bytes(), &response))
				assert.Equal(t, tt.failed, response.Failed)
			}
		})
	}
}
