package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type createAdminRequest struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestHandler_CreateAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		args         createAdminRequest
		mockSetup    func(service *MockadminService)
		expectedCode int
	}{
		{
			name: "Missing fields",
			args: createAdminRequest{
				ID:       1,
				Username: "admin",
			},
			mockSetup:    func(service *MockadminService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name: "Correct request",
			args: createAdminRequest{
				ID:       1,
				Username: "admin",
				Password: "password",
			},
			mockSetup: func(adminService *MockadminService) {
				adminService.EXPECT().CreateAdmin(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name: "Such id exists",
			args: createAdminRequest{
				ID:       1,
				Username: "dasdasd",
				Password: "password",
			},
			mockSetup: func(adminService *MockadminService) {
				adminService.EXPECT().CreateAdmin(gomock.Any(), gomock.Any()).Return(admin.ErrIDUsed).Times(1)
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name: "Such username exists",
			args: createAdminRequest{
				ID:       123,
				Username: "admin",
				Password: "password",
			},
			mockSetup: func(adminService *MockadminService) {
				adminService.EXPECT().CreateAdmin(gomock.Any(), gomock.Any()).Return(admin.ErrUsernameUsed).Times(1)
			},
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := NewMockadminService(ctrl)
			tt.mockSetup(mockService)

			reqBody, err := json.Marshal(tt.args)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/admins", bytes.NewReader(reqBody))
			res := httptest.NewRecorder()
			handler := NewHandler(mockService)

			handler.CreateAdmin(t.Context(), res, req)

			assert.Equal(t, tt.expectedCode, res.Code)
		})
	}
}
