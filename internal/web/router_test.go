package web

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Rhymond/go-money"
	"github.com/jackc/pgx/v4"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"gitlab.ozon.dev/alexplay1224/homework/internal/config"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

type request struct {
	method string
	path   string
	body   []byte
}

func TestApp_Run(t *testing.T) {
	t.Parallel()

	password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	tests := []struct {
		name         string
		args         request
		authorized   bool
		mockSetup    func(MockorderStorage, MockadminStorage, MockauditLoggerStorage, MocktxManager)
		expectedCode int
	}{
		{
			name: "valid get orders",
			args: request{
				method: http.MethodGet,
				path:   "/orders",
			},
			authorized: true,
			mockSetup: func(mockOrderStorage MockorderStorage, mockAdminStorage MockadminStorage,
				_ MockauditLoggerStorage, _ MocktxManager) {
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockOrderStorage.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Order{}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "not valid get orders",
			args: request{
				method: http.MethodGet,
				path:   "/order",
			},
			authorized: true,
			mockSetup: func(_ MockorderStorage, _ MockadminStorage,
				_ MockauditLoggerStorage, _ MocktxManager) {
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "not authorised get orders",
			args: request{
				method: http.MethodGet,
				path:   "/orders",
			},
			authorized: false,
			mockSetup: func(_ MockorderStorage, _ MockadminStorage,
				_ MockauditLoggerStorage, _ MocktxManager) {
			},
			expectedCode: http.StatusUnauthorized,
		},
		{
			name: "valid post orders",
			args: request{
				method: http.MethodPost,
				path:   "/orders",
				body: []byte(`{"id": 111111111111,"user_id":52,"weight":100,"price":{"amount":1000000,"currency":"RUB"},
								"packaging":2,"extra_packaging":3,"expiry_date":"4025-03-10T00:00:00Z"}`),
			},
			authorized: true,
			mockSetup: func(mockOrderStorage MockorderStorage, mockAdminStorage MockadminStorage,
				_ MockauditLoggerStorage, tx MocktxManager) {
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				tx.EXPECT().RunRepeatableRead(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) error {
						return f(ctx, nil)
					}).Return(nil)
				mockOrderStorage.EXPECT().AddOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockOrderStorage.EXPECT().Contains(gomock.Any(), gomock.Any(), gomock.Any()).Return(false, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "not valid post orders",
			args: request{
				method: http.MethodPost,
				path:   "/order",
				body: []byte(`{"id": 1009,"user_id":52,"weight":100,"price":{"amount":1000000,"currency":"RUB"},
								"packaging":2,"extra_packaging":3,"expiry_date":"2025-03-10T00:00:00Z"}`),
			},
			authorized: true,
			mockSetup: func(_ MockorderStorage, _ MockadminStorage, _ MockauditLoggerStorage, _ MocktxManager) {
			},
			expectedCode: http.StatusNotFound,
		},
		{
			name: "valid delete orders",
			args: request{
				method: http.MethodDelete,
				path:   "/orders/123",
			},
			authorized: true,
			mockSetup: func(mockOrderStorage MockorderStorage, mockAdminStorage MockadminStorage,
				_ MockauditLoggerStorage, mocktxManager MocktxManager) {
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mocktxManager.EXPECT().RunSerializable(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) error {
						return f(ctx, nil)
					})
				mockOrderStorage.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Order{}, nil)
				mockOrderStorage.EXPECT().Contains(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockOrderStorage.EXPECT().RemoveOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid process orders",
			args: request{
				method: http.MethodPost,
				path:   "/orders/process",
				body:   []byte(`{"id":4,"user_id":789,"action":"give"}`),
			},
			authorized: true,
			mockSetup: func(mockOrderStorage MockorderStorage, mockAdminStorage MockadminStorage,
				logger MockauditLoggerStorage, tx MocktxManager) {
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				tx.EXPECT().RunSerializable(gomock.Any(), gomock.Any()).Return(nil).
					DoAndReturn(func(ctx context.Context, f func(ctx context.Context, tx pgx.Tx) error) error {
						return f(ctx, nil)
					})
				logger.EXPECT().GetAndMarkLogs(gomock.Any(), gomock.Any()).Return(nil, nil)
				mockOrderStorage.EXPECT().UpdateOrder(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockOrderStorage.EXPECT().Contains(gomock.Any(), gomock.Any(), gomock.Any()).Return(true, nil)
				mockOrderStorage.EXPECT().GetByID(gomock.Any(), gomock.Any(), gomock.Any()).Return(models.Order{
					ID: 4, UserID: 789, Weight: 1233, Price: *money.New(22222, money.RUB),
					Status: 1, ExpiryDate: time.Now().Add(time.Hour)}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid post admins",
			args: request{
				method: http.MethodPost,
				path:   "/admins",
				body:   []byte(`{"id":52,"username":"sdasds","password":"give"}`),
			},
			authorized: false,
			mockSetup: func(_ MockorderStorage, mockAdminStorage MockadminStorage, _ MockauditLoggerStorage, _ MocktxManager) {
				mockAdminStorage.EXPECT().CreateAdmin(gomock.Any(), gomock.Any()).Return(nil)
				mockAdminStorage.EXPECT().ContainsID(gomock.Any(), gomock.Any()).Return(false, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(false, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid delete admins",
			args: request{
				method: http.MethodDelete,
				path:   "/admins/asdasd",
				body:   []byte(`{"password":"password"}`),
			},
			authorized: false,
			mockSetup: func(_ MockorderStorage, mockAdminStorage MockadminStorage,
				_ MockauditLoggerStorage, _ MocktxManager) {
				mockAdminStorage.EXPECT().DeleteAdmin(gomock.Any(), gomock.Any()).Return(nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "asdasd", Password: string(password)}, nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	ctrl := gomock.NewController(t)

	t.Cleanup(func() {
		ctrl.Finish()
	})

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockTxManager := NewMocktxManager(ctrl)

			mockOrderStorage := NewMockorderStorage(ctrl)
			mockAdminStorage := NewMockadminStorage(ctrl)
			mockLogStorage := NewMockauditLoggerStorage(ctrl)
			app, _ := NewApp(context.Background(), config.Config{}, mockOrderStorage, mockAdminStorage, mockLogStorage,
				mockTxManager, 2, 5, 500*time.Second)
			app.SetupRoutes(context.Background())

			tt.mockSetup(*mockOrderStorage, *mockAdminStorage, *mockLogStorage, *mockTxManager)

			var authHeader string
			req, err := http.NewRequestWithContext(context.Background(), tt.args.method, tt.args.path,
				bytes.NewReader(tt.args.body))
			require.NoError(t, err)
			if tt.authorized {
				username := "user"
				password := "password"
				auth := username + ":" + password
				authHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
				req.Header.Set("Authorization", authHeader)
			}

			if tt.name == "valid post orders" {
				fmt.Print(1)
			}

			res := httptest.NewRecorder()
			app.Router.ServeHTTP(res, req)

			require.Equal(t, tt.expectedCode, res.Code)
		})
	}
}
