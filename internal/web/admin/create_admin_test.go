//go:build unit

package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.ozon.dev/alexplay1224/homework/internal/mocks/service"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

type createAdminRequest struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func TestHandler_CreateAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		args      createAdminRequest
		mockSetup func(service *service.MockadminService)
		want      int
	}{
		{
			name: "Missing fields",
			args: createAdminRequest{
				ID:       1,
				Username: "admin",
			},
			mockSetup: func(service *service.MockadminService) {},
			want:      http.StatusBadRequest,
		},
		{
			name: "Correct request",
			args: createAdminRequest{
				ID:       1,
				Username: "admin",
				Password: "password",
			},
			mockSetup: func(adminService *service.MockadminService) {
				adminService.EXPECT().CreateAdmin(gomock.Any(), gomock.Any()).Return(nil).Times(1)
			},
			want: http.StatusOK,
		},
		{
			name: "Such id exists",
			args: createAdminRequest{
				ID:       1,
				Username: "dasdasd",
				Password: "password",
			},
			mockSetup: func(adminService *service.MockadminService) {
				adminService.EXPECT().CreateAdmin(gomock.Any(), gomock.Any()).Return(admin.ErrIDUsed).Times(1)
			},
			want: http.StatusInternalServerError,
		},
		{
			name: "Such username exists",
			args: createAdminRequest{
				ID:       123,
				Username: "admin",
				Password: "password",
			},
			mockSetup: func(adminService *service.MockadminService) {
				adminService.EXPECT().CreateAdmin(gomock.Any(), gomock.Any()).Return(admin.ErrUsernameUsed).Times(1)
			},
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockService := service.NewMockadminService(ctrl)
			tt.mockSetup(mockService)

			reqBody, err := json.Marshal(tt.args)
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(http.MethodPost, "/admins", bytes.NewReader(reqBody))
			res := httptest.NewRecorder()
			handler := NewHandler(mockService)

			handler.CreateAdmin(t.Context(), res, req)

			assert.Equal(t, tt.want, res.Code)
		})
	}
}
