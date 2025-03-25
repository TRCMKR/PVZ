//go:build integration

package order

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	order_Service "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	order_Handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/order"
	"gitlab.ozon.dev/alexplay1224/homework/tests/integration"
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
				"user_id": "52",
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
				"count": "4",
			},
			expectedStatus: http.StatusOK,
			expectedCount:  4,
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
	orderService := order_Service.NewService(ordersRepo)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
		defer db.Close()
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/orders", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			res := httptest.NewRecorder()
			handler := order_Handler.NewHandler(orderService)

			handler.GetOrders(ctx, res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)

			if tt.expectedStatus == http.StatusOK {
				var response struct {
					Count  int            `json:"count"`
					Orders []models.Order `json:"orders"`
				}
				err := json.NewDecoder(res.Body).Decode(&response)
				require.NoError(t, err)
				assert.Equal(t, tt.expectedCount, response.Count)
			}
		})
	}
}
