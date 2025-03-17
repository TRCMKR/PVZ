package order

import (
	"context"
	"encoding/json"
	"github.com/Rhymond/go-money"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
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

func TestOrderHandler_GetOrders(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "Valid request with filters",
			queryParams: map[string]string{
				"user_id": "123",
				"count":   "10",
				"page":    "0",
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "Invalid count",
			queryParams: map[string]string{
				"user_id": "123",
				"count":   "invalid",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name: "Invalid filter (invalid weight)",
			queryParams: map[string]string{
				"weight": "not-a-number",
			},
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name: "No filters, page and count provided",
			queryParams: map[string]string{
				"count": "5",
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "Valid request with filters",
			queryParams: map[string]string{
				"order_id": "1",
				"user_id":  "123",
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
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

			req := httptest.NewRequest(http.MethodGet, "/orders", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			res := httptest.NewRecorder()
			handler := orderHandlerPkg.NewHandler(orderService)

			handler.GetOrders(ctx, res, req)

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
			}
		})
	}
}
