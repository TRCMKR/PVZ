//go:build integration

package order

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	order_Service "gitlab.ozon.dev/alexplay1224/homework/internal/service/order"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	order_Handler "gitlab.ozon.dev/alexplay1224/homework/internal/web/order"
	"gitlab.ozon.dev/alexplay1224/homework/tests/integration"
)

func TestOrderHandler_DeleteOrder(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		orderID        int
		expectedStatus int
	}{
		{
			name:           "Delete existing order",
			orderID:        1,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Delete non-existing order",
			orderID:        9999,
			expectedStatus: http.StatusInternalServerError,
		},
	}

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

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/orders/%d", tt.orderID), nil)
			req = mux.SetURLVars(req, map[string]string{
				order_Handler.OrderIDParam: strconv.Itoa(tt.orderID),
			})
			res := httptest.NewRecorder()
			handler := order_Handler.NewHandler(orderService)

			handler.DeleteOrder(t.Context(), res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)
			if tt.expectedStatus == http.StatusOK {
				require.Equal(t, "success", res.Body.String())
			}
		})
	}
}
