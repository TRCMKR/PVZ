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
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type updareAdminRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

func TestUpdateAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		username  string
		args      updareAdminRequest
		mockSetup func(service *service.MockadminService)
		want      int
	}{
		{
			name:     "Missing fields",
			username: "test",
			args: updareAdminRequest{
				Password: "123123",
			},
			mockSetup: func(service *service.MockadminService) {},
			want:      http.StatusBadRequest,
		},
		{
			name:     "User not found",
			username: "test",
			args: updareAdminRequest{
				Password:    "password",
				NewPassword: "new_password",
			},
			mockSetup: func(adminService *service.MockadminService) {
				adminService.EXPECT().UpdateAdmin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(admin.ErrAdminDoesntExist).Times(1)
			},
			want: http.StatusInternalServerError,
		},
		{
			name:     "Correct request",
			username: "test",
			args: updareAdminRequest{
				Password:    "password",
				NewPassword: "new_password",
			},
			mockSetup: func(adminService *service.MockadminService) {
				adminService.EXPECT().UpdateAdmin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
			},
			want: http.StatusOK,
		},
		{
			name:     "Incorrect password",
			username: "test",
			args: updareAdminRequest{
				Password:    "password",
				NewPassword: "new_password",
			},
			mockSetup: func(adminService *service.MockadminService) {
				adminService.EXPECT().UpdateAdmin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(admin.ErrIDUsed).Times(1)
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

			req := httptest.NewRequest(http.MethodPost, "/admin/"+tt.username, bytes.NewBuffer(reqBody))
			req = mux.SetURLVars(req, map[string]string{
				AdminUsernameParam: tt.username,
			})
			res := httptest.NewRecorder()

			handler := NewHandler(mockService)
			handler.UpdateAdmin(t.Context(), res, req)

			assert.Equal(t, tt.want, res.Code)
		})
	}
}
