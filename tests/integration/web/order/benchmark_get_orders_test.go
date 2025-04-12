package order

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/facade"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/tx_manager"
	web "gitlab.ozon.dev/alexplay1224/homework/internal/web/http"
	"gitlab.ozon.dev/alexplay1224/homework/tests/integration"
)

func setup() (string, string, error) {
	ctx := context.Background()

	rootDir, _ := config.GetRootDir()
	err := config.InitEnv(rootDir + "/.env.test")
	if err != nil {
		return "", "", err
	}
	cfg := config.NewConfig()

	connStr, _, _ := integration.InitPostgresContainer(ctx, cfg)
	url := "/orders"

	return connStr, url, nil
}

func sendGetOrdersRequest(ctx context.Context, wg *sync.WaitGroup, url string, username string, password string) {
	defer wg.Done()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Printf("failed to create request: %v", err)

		return
	}

	auth := username + ":" + password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("failed to execute request: %v", err)

		return
	}
	defer resp.Body.Close()

	// Проверяем статус-код ответа
	if resp.StatusCode != http.StatusOK {
		log.Printf("unexpected status code: %d", resp.StatusCode)
	}
}

func BenchmarkOrderHandler_GetOrders_Cache(b *testing.B) {
	ctx := context.Background()
	connStr, url, err := setup()
	require.NoError(b, err)

	db, _ := postgres.NewDB(ctx, connStr)

	txManager := tx_manager.NewTxManager(db)

	logger := zap.NewNop()

	ordersRepo := repository.NewOrdersRepo(logger, db)
	ordersFacade := facade.NewOrderFacade(ctx, ordersRepo, 10000)

	adminsRepo := repository.NewAdminsRepo(logger, db)
	adminsFacade := facade.NewAdminFacade(adminsRepo, 10000)

	logsRepo := repository.NewLogsRepo(db)

	app, _ := web.NewApp(ctx, config.Config{}, logger, ordersFacade, adminsFacade, logsRepo, txManager,
		2, 5, 500*time.Millisecond)
	app.SetupRoutes(ctx)

	server := httptest.NewServer(app.Router)
	go func() {
		<-ctx.Done()
		server.Close()
	}()

	numRequests := b.N

	b.ResetTimer()

	maxConcurrentRequests := 10
	sem := make(chan struct{}, maxConcurrentRequests)
	var wg sync.WaitGroup
	wg.Add(numRequests)
	for i := 0; i < b.N; i++ {
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			sendGetOrdersRequest(ctx, &wg, server.URL+url, "test", "12345678")
		}()
	}
	wg.Wait()
	db.Close()
}

func BenchmarkOrderHandler_GetOrders_NoAdminCache(b *testing.B) {
	ctx := context.Background()
	connStr, url, err := setup()
	require.NoError(b, err)

	db, _ := postgres.NewDB(ctx, connStr)

	txManager := tx_manager.NewTxManager(db)

	logger := zap.NewNop()

	ordersRepo := repository.NewOrdersRepo(logger, db)
	ordersFacade := facade.NewOrderFacade(ctx, ordersRepo, 10000)

	adminsRepo := repository.NewAdminsRepo(logger, db)

	logsRepo := repository.NewLogsRepo(db)

	app, _ := web.NewApp(ctx, config.Config{}, logger, ordersFacade, adminsRepo, logsRepo, txManager,
		2, 5, 500*time.Millisecond)
	app.SetupRoutes(ctx)

	server := httptest.NewServer(app.Router)
	go func() {
		<-ctx.Done()
		server.Close()
	}()

	numRequests := b.N

	b.ResetTimer()

	maxConcurrentRequests := 10
	sem := make(chan struct{}, maxConcurrentRequests)
	var wg sync.WaitGroup
	wg.Add(numRequests)
	for i := 0; i < b.N; i++ {
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			sendGetOrdersRequest(ctx, &wg, server.URL+url, "test", "12345678")
		}()
	}
	wg.Wait()
	db.Close()
}

func BenchmarkOrderHandler_GetOrders_NoOrderCache(b *testing.B) {
	ctx := context.Background()
	connStr, url, err := setup()
	require.NoError(b, err)

	db, _ := postgres.NewDB(ctx, connStr)

	logger := zap.NewNop()

	txManager := tx_manager.NewTxManager(db)

	ordersRepo := repository.NewOrdersRepo(logger, db)
	adminsRepo := repository.NewAdminsRepo(logger, db)

	adminsFacade := facade.NewAdminFacade(adminsRepo, 10000)

	logsRepo := repository.NewLogsRepo(db)

	app, _ := web.NewApp(ctx, config.Config{}, logger, ordersRepo, adminsFacade, logsRepo, txManager,
		2, 5, 500*time.Millisecond)
	app.SetupRoutes(ctx)

	server := httptest.NewServer(app.Router)
	go func() {
		<-ctx.Done()
		server.Close()
	}()

	numRequests := b.N

	b.ResetTimer()

	maxConcurrentRequests := 10
	sem := make(chan struct{}, maxConcurrentRequests)
	var wg sync.WaitGroup
	wg.Add(numRequests)
	for i := 0; i < b.N; i++ {
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			sendGetOrdersRequest(ctx, &wg, server.URL+url, "test", "12345678")
		}()
	}
	wg.Wait()
	db.Close()
}

func BenchmarkOrderHandler_GetOrders_NoCache(b *testing.B) {
	ctx := context.Background()
	connStr, url, err := setup()
	require.NoError(b, err)

	db, _ := postgres.NewDB(ctx, connStr)

	txManager := tx_manager.NewTxManager(db)

	logger := zap.NewNop()

	ordersRepo := repository.NewOrdersRepo(logger, db)

	adminsRepo := repository.NewAdminsRepo(logger, db)

	logsRepo := repository.NewLogsRepo(db)

	app, _ := web.NewApp(ctx, config.Config{}, logger, ordersRepo, adminsRepo, logsRepo, txManager,
		2, 5, 500*time.Millisecond)
	app.SetupRoutes(ctx)

	server := httptest.NewServer(app.Router)
	go func() {
		<-ctx.Done()
		server.Close()
	}()

	numRequests := b.N

	b.ResetTimer()

	maxConcurrentRequests := 10
	sem := make(chan struct{}, maxConcurrentRequests)
	var wg sync.WaitGroup
	wg.Add(numRequests)
	for i := 0; i < b.N; i++ {
		sem <- struct{}{}
		go func() {
			defer func() { <-sem }()
			sendGetOrdersRequest(ctx, &wg, server.URL+url, "test", "12345678")
		}()
	}
	wg.Wait()
	db.Close()
}
