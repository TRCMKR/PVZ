//go:build integration

package order

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/require"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/facade"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/repository"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres/tx_manager"
	"gitlab.ozon.dev/alexplay1224/homework/internal/web"
	"gitlab.ozon.dev/alexplay1224/homework/tests/integration"
)

func sendRequest(ctx context.Context, wg *sync.WaitGroup, url string, requestData createOrderRequest,
	username string, password string, ch chan<- error) {
	defer wg.Done()

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		ch <- fmt.Errorf("error marshalling request data: %w", err)

		return
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonData))
	if err != nil {
		ch <- fmt.Errorf("error creating request: %w", err)

		return
	}

	req.Header.Set("Content-Type", "application/json")
	auth := username + ":" + password
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(auth))
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		ch <- fmt.Errorf("error making request: %w", err)

		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- fmt.Errorf("expected 200 OK, got %v", resp.StatusCode)

		return
	}

	ch <- nil
}

func generateOrderRequests(count int) []createOrderRequest {
	var requests []createOrderRequest

	for i := 1; i <= count; i++ {
		request := createOrderRequest{
			ID:             i + 1000,
			UserID:         789,
			Weight:         100,
			Price:          *money.New(1000, money.RUB),
			Packaging:      1,
			ExtraPackaging: 0,
			ExpiryDate:     time.Now().Add(time.Hour * 24),
		}
		requests = append(requests, request)
	}

	return requests
}

func TestOrderHandlerRps_CreateOrder(t *testing.T) {
	t.Parallel()

	ctx := context.Background()
	rootDir, err := config.GetRootDir()
	require.NoError(t, err)
	config.InitEnv(rootDir + "/.env.test")
	cfg := *config.NewConfig()

	connStr, pgContainer, err := integration.InitPostgresContainer(ctx, cfg)
	require.NoError(t, err)
	db, err := postgres.NewDB(ctx, connStr)
	require.NoError(t, err)

	txManager := tx_manager.NewTxManager(db)

	ordersRepo := repository.NewOrderRepo(txManager)
	ordersFacade := facade.NewOrderFacade(ctx, ordersRepo, 10000)
	adminsRepo := repository.NewAdminRepo(db)
	adminsFacade := facade.NewAdminFacade(adminsRepo, 10000)
	logsRepo := repository.NewLogsRepo(db)

	app, _ := web.NewApp(ctx, ordersFacade, adminsFacade, logsRepo, txManager, 2, 5, 500*time.Millisecond)
	app.SetupRoutes(ctx)

	server := httptest.NewServer(app.Router)
	go func() {
		<-ctx.Done()
		server.Close()
	}()

	t.Cleanup(func() {
		if err = pgContainer.Terminate(context.Background()); err != nil {
			t.Fatalf("failed to terminate pgContainer: %s", err)
		}
		db.Close()
	})

	numRequests := 60
	requests := generateOrderRequests(numRequests)

	var wg sync.WaitGroup
	errChan := make(chan error, numRequests)

	startTime := time.Now()
	wg.Add(numRequests)
	maxConcurrentRequests := 10
	sem := make(chan struct{}, maxConcurrentRequests)

	for i := 0; i < numRequests; i++ {
		sem <- struct{}{}
		go func(i int) {
			defer func() { <-sem }()
			sendRequest(ctx, &wg, server.URL+"/orders", requests[i], "test", "12345678", errChan)
		}(i)
	}
	wg.Wait()
	close(errChan)

	duration := time.Since(startTime)
	rps := float64(numRequests)
	t.Logf("Requests per second: %.2f, time: %s", rps, duration)

	select {
	case <-ctx.Done():
		t.Fatal("timed out waiting for orders to be created")
	default:
	}

	for err = range errChan {
		require.NoError(t, err)
	}

	require.GreaterOrEqual(t, rps, 10.0, "RPS should be greater or equal to 10")
}
