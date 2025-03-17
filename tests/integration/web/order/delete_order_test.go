package order

import (
	"context"
	"fmt"
	"github.com/Rhymond/go-money"
	"github.com/gorilla/mux"
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
	"strconv"
	"testing"
	"time"
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
			ExpiryDate:  time.Now().AddDate(-1, 0, 0),
		},
	}

	for _, order := range orders {
		err = ordersRepo.AddOrder(ctx, order)
		if err != nil {
			t.Fatal(err)
		}
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

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/orders/%d", tt.orderID), nil)
			req = mux.SetURLVars(req, map[string]string{
				orderHandlerPkg.OrderIDParam: strconv.Itoa(tt.orderID),
			})
			res := httptest.NewRecorder()
			handler := orderHandlerPkg.NewHandler(orderService)

			handler.DeleteOrder(ctx, res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)
		})
	}
}
