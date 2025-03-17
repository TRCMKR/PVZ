package web

import (
	"bytes"
	"encoding/base64"
	"github.com/stretchr/testify/require"
	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"testing"
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
		mockSetup    func(mockOrderService MockorderStorage, mockAdminStorage MockadminStorage)
		expectedCode int
	}{
		{
			name: "valid get orders",
			args: request{
				method: http.MethodGet,
				path:   "/orders",
			},
			authorized: true,
			mockSetup: func(mockOrderService MockorderStorage, mockAdminStorage MockadminStorage) {
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockOrderService.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
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
			mockSetup: func(mockOrderService MockorderStorage, mockAdminStorage MockadminStorage) {
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
			mockSetup: func(mockOrderService MockorderStorage, mockAdminStorage MockadminStorage) {
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockOrderService.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Order{}, nil)
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
			mockSetup: func(mockOrderService MockorderStorage, mockAdminStorage MockadminStorage) {
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockOrderService.EXPECT().GetOrders(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.Order{}, nil)
				mockOrderService.EXPECT().AddOrder(gomock.Any(), gomock.Any()).Return(nil)
				mockOrderService.EXPECT().Contains(gomock.Any(), gomock.Any()).Return(false, nil).Times(2)
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
			mockSetup: func(mockOrderService MockorderStorage, mockAdminStorage MockadminStorage) {
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
			mockSetup: func(mockOrderService MockorderStorage, mockAdminStorage MockadminStorage) {
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockOrderService.EXPECT().RemoveOrder(gomock.Any(), gomock.Any()).Return(nil)
				mockOrderService.EXPECT().Contains(gomock.Any(), gomock.Any()).Return(true, nil)
				mockOrderService.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(models.Order{}, nil)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "valid process orders",
			args: request{
				method: http.MethodPost,
				path:   "/orders/process",
				body:   []byte(`{"user_id":52,"order_ids":[1],"action":"give"}`),
			},
			authorized: true,
			mockSetup: func(mockOrderService MockorderStorage, mockAdminStorage MockadminStorage) {
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "user", Password: string(password)}, nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockOrderService.EXPECT().UpdateOrder(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)
				mockOrderService.EXPECT().Contains(gomock.Any(), gomock.Any()).Return(true, nil)
				mockOrderService.EXPECT().GetByID(gomock.Any(), gomock.Any()).Return(models.Order{}, nil)
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
			mockSetup: func(mockOrderService MockorderStorage, mockAdminStorage MockadminStorage) {
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
			mockSetup: func(mockOrderService MockorderStorage, mockAdminStorage MockadminStorage) {
				mockAdminStorage.EXPECT().DeleteAdmin(gomock.Any(), gomock.Any()).Return(nil)
				mockAdminStorage.EXPECT().ContainsUsername(gomock.Any(), gomock.Any()).Return(true, nil)
				mockAdminStorage.EXPECT().GetAdminByUsername(gomock.Any(), gomock.Any()).
					Return(models.Admin{ID: 0, Username: "asdasd", Password: string(password)}, nil)
			},
			expectedCode: http.StatusOK,
		},
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mockOrderStorage := NewMockorderStorage(ctrl)
			mockAdminStorage := NewMockadminStorage(ctrl)
			app := NewApp(mockOrderStorage, mockAdminStorage)
			app.SetupRoutes(t.Context())

			tt.mockSetup(*mockOrderStorage, *mockAdminStorage)

			var authHeader string
			req, err := http.NewRequest(tt.args.method, tt.args.path, bytes.NewReader(tt.args.body))
			require.NoError(t, err)
			if tt.authorized {
				username := "user"
				password := "password"
				auth := username + ":" + password
				authHeader = "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
				req.Header.Set("Authorization", authHeader)
			}

			res := httptest.NewRecorder()
			app.router.ServeHTTP(res, req)

			require.Equal(t, tt.expectedCode, res.Code)
		})
	}
}
