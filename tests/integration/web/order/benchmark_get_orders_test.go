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

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/facade"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/tx_manager"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"
	"gitlab.ozon.dev/alexplay1224/homework/tests/integration"
)

func setup() (string, string) {
	ctx := context.Background()

	rootDir, _ := config.GetRootDir()
	config.InitEnv(rootDir + "/.env.test")
	cfg := config.NewConfig()

	connStr, _, _ := integration.InitPostgresContainer(ctx, cfg)
	url := "/orders"

	return connStr, url
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
	connStr, url := setup()
	db, _ := postgres.NewDB(ctx, connStr)

	txManager := tx_manager.NewTxManager(db)
	ordersRepo := repository.NewOrdersRepo(db)
	ordersFacade := facade.NewOrderFacade(ctx, ordersRepo, 10000)
	adminsRepo := repository.NewAdminsRepo(db)
	adminsFacade := facade.NewAdminFacade(adminsRepo, 10000)
	logsRepo := repository.NewLogsRepo(db)

	app, _ := web.NewApp(ctx, config.Config{}, ordersFacade, adminsFacade, logsRepo, txManager, 2, 5, 500*time.Millisecond)
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
	connStr, url := setup()
	db, _ := postgres.NewDB(ctx, connStr)

	txManager := tx_manager.NewTxManager(db)
	ordersRepo := repository.NewOrdersRepo(db)
	ordersFacade := facade.NewOrderFacade(ctx, ordersRepo, 10000)
	adminsRepo := repository.NewAdminsRepo(db)
	logsRepo := repository.NewLogsRepo(db)

	app, _ := web.NewApp(ctx, config.Config{}, ordersFacade, adminsRepo, logsRepo, txManager, 2, 5, 500*time.Millisecond)
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
	connStr, url := setup()
	db, _ := postgres.NewDB(ctx, connStr)

	txManager := tx_manager.NewTxManager(db)
	ordersRepo := repository.NewOrdersRepo(db)
	adminsRepo := repository.NewAdminsRepo(db)
	adminsFacade := facade.NewAdminFacade(adminsRepo, 10000)
	logsRepo := repository.NewLogsRepo(db)

	app, _ := web.NewApp(ctx, config.Config{}, ordersRepo, adminsFacade, logsRepo, txManager, 2, 5, 500*time.Millisecond)
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
	connStr, url := setup()
	db, _ := postgres.NewDB(ctx, connStr)

	txManager := tx_manager.NewTxManager(db)
	ordersRepo := repository.NewOrdersRepo(db)
	adminsRepo := repository.NewAdminsRepo(db)
	logsRepo := repository.NewLogsRepo(db)

	app, _ := web.NewApp(ctx, config.Config{}, ordersRepo, adminsRepo, logsRepo, txManager, 2, 5, 500*time.Millisecond)
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
