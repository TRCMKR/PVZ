//go:build integration

package order

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	order_Handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/http/order"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	order_Service "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/tx_manager"
	"gitlab.ozon.dev/alexplay1224/homework/tests/integration"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrderHandler_UpdateOrders(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Update existing order",
			requestBody: `{
                "user_id": 789,
                "id": 4,
                "action": "give"
            }`,
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name: "Update non-existing order",
			requestBody: `{
                "user_id": 123,
                "id": 1232,
                "action": "give"
            }`,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	ctx := context.Background()
	rootDir, err := config.GetRootDir()
	require.NoError(t, err)
	config.InitEnv(rootDir + "/.env.test")
	cfg := config.NewConfig()

	connStr, pgContainer, err := integration.InitPostgresContainer(t.Context(), cfg)
	require.NoError(t, err)
	db, err := postgres.NewDB(t.Context(), connStr)
	require.NoError(t, err)

	txManager := tx_manager.NewTxManager(db)
	ordersRepo := repository.NewOrdersRepo(db)
	orderService := order_Service.NewService(ordersRepo, txManager)

	t.Cleanup(func() {
		if err := pgContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
		defer db.Close()
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			reqBody := []byte(tt.requestBody)

			req := httptest.NewRequest(http.MethodPost, "/orders/process", bytes.NewReader(reqBody))
			res := httptest.NewRecorder()
			handler := order_Handler.NewHandler(orderService)

			handler.UpdateOrder(ctx, res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedBody, res.Body.String())
			}
		})
	}
}
