//go:build unit

package admin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.ozon.dev/alexplay1224/homework/internal/mocks/service"
	"gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

type adminDeleteRequest struct {
	Password string `json:"password"`
}

func TestHandler_DeleteAdmin(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		username  string
		args      adminDeleteRequest
		mockSetup func(service *service.MockadminService)
		want      int
	}{
		{
			name:     "Correct deletion",
			username: "admin",
			args: adminDeleteRequest{
				Password: "12345678",
			},
			mockSetup: func(adminService *service.MockadminService) {
				adminService.EXPECT().DeleteAdmin(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
			},
			want: http.StatusOK,
		},
		{
			name:     "Wrong password",
			username: "admin",
			args: adminDeleteRequest{
				Password: "12345678",
			},
			mockSetup: func(adminService *service.MockadminService) {
				adminService.EXPECT().DeleteAdmin(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(admin.ErrWrongPassword).Times(1)
			},
			want: http.StatusInternalServerError,
		},
		{
			name:     "Invalid password",
			username: "admin",
			args: adminDeleteRequest{
				Password: "",
			},
			mockSetup: func(adminService *service.MockadminService) {},
			want:      http.StatusBadRequest,
		},
		{
			name:     "No such admin",
			username: "fake_admin",
			args: adminDeleteRequest{
				Password: "12345678",
			},
			mockSetup: func(adminService *service.MockadminService) {
				adminService.EXPECT().DeleteAdmin(gomock.Any(), gomock.Any(), gomock.Any()).
					Return(admin.ErrAdminDoesntExist).Times(1)
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

			req := httptest.NewRequest(http.MethodDelete, "/admin/"+tt.username, bytes.NewReader(reqBody))
			req = mux.SetURLVars(req, map[string]string{
				AdminUsernameParam: tt.username,
			})
			res := httptest.NewRecorder()

			if tt.name == "Invalid password" {
				fmt.Print(1)
			}

			handler := NewHandler(mockService)
			handler.DeleteAdmin(t.Context(), res, req)

			assert.Equal(t, tt.want, res.Code)
		})
	}
}
