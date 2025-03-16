//go:build integration

package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"

	"github.com/Rhymond/go-money"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	pgcontainer "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

func initPostgresContainer(ctx context.Context, cfg config.Config) (string, *pgcontainer.PostgresContainer, error) {
	pgContainer, err := pgcontainer.Run(ctx, "postgres:14-alpine",
		pgcontainer.WithDatabase(cfg.DBName()),
		pgcontainer.WithUsername(cfg.Username()),
		pgcontainer.WithPassword(cfg.Password()),
		testcontainers.WithLogger(log.Default()),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(5*time.Second)),
	)
	if err != nil {
		return "", nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return "", nil, err
	}
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return "", nil, err
	}
	defer db.Close()

	time.Sleep(2 * time.Second)

	goose.SetDialect("postgres")
	if err := goose.Up(db, "../../migrations"); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return connStr, pgContainer, nil
}

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
	config.InitEnv("../../.env.test")
	cfg := *config.NewConfig()

	connStr, pgContainer, err := initPostgresContainer(ctx, cfg)
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

			orderService := service.OrderService{
				Storage: ordersRepo,
			}

			handler := web.OrderHandler{
				OrderService: &orderService,
			}

			if tt.name == "Correct order" {
				fmt.Print(1)
			}

			handler.CreateOrder(ctx, res, req)

			assert.Equal(t, tt.want, res.Code)
		})
	}
}

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
	config.InitEnv("../../.env.test")
	cfg := *config.NewConfig()

	connStr, pgContainer, err := initPostgresContainer(ctx, cfg)
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
			orderService := service.OrderService{
				Storage: ordersRepo,
			}

			req := httptest.NewRequest(http.MethodGet, "/orders", nil)
			q := req.URL.Query()
			for key, value := range tt.queryParams {
				q.Add(key, value)
			}
			req.URL.RawQuery = q.Encode()

			res := httptest.NewRecorder()
			handler := web.OrderHandler{
				OrderService: &orderService,
			}

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
	config.InitEnv("../../.env.test")
	cfg := *config.NewConfig()

	connStr, pgContainer, err := initPostgresContainer(ctx, cfg)
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
			orderService := service.OrderService{
				Storage: ordersRepo,
			}

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/orders/%d", tt.orderID), nil)
			req = mux.SetURLVars(req, map[string]string{
				web.OrderIDParam: strconv.Itoa(tt.orderID),
			})
			res := httptest.NewRecorder()
			handler := web.OrderHandler{
				OrderService: &orderService,
			}

			handler.DeleteOrder(ctx, res, req)

			assert.Equal(t, tt.expectedStatus, res.Code)
		})
	}
}

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
	config.InitEnv("../../.env.test")
	cfg := *config.NewConfig()

	connStr, pgContainer, err := initPostgresContainer(ctx, cfg)
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
			orderService := service.OrderService{
				Storage: ordersRepo,
			}

			reqBody := []byte(tt.requestBody)

			req := httptest.NewRequest(http.MethodPost, "/orders/process", bytes.NewReader(reqBody))
			res := httptest.NewRecorder()
			handler := web.OrderHandler{
				OrderService: &orderService,
			}

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
