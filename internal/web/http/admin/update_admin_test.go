package admin

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"gitlab.ozon.dev/alexplay1224/homework/internal/service/admin"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type updareAdminRequest struct {
	Password    string `json:"password"`
	NewPassword string `json:"new_password"`
}

func TestUpdateAdmin(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name         string
		username     string
		args         updareAdminRequest
		mockSetup    func(service *MockadminService)
		expectedCode int
	}{
		{
			name:     "Missing fields",
			username: "test",
			args: updareAdminRequest{
				Password: "123123",
			},
			mockSetup:    func(_ *MockadminService) {},
			expectedCode: http.StatusBadRequest,
		},
		{
			name:     "User not found",
			username: "test",
			args: updareAdminRequest{
				Password:    "password",
				NewPassword: "new_password",
			},
			mockSetup: func(adminService *MockadminService) {
				adminService.EXPECT().UpdateAdmin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(admin.ErrAdminDoesntExist).Times(1)
			},
			expectedCode: http.StatusInternalServerError,
		},
		{
			name:     "Correct request",
			username: "test",
			args: updareAdminRequest{
				Password:    "password",
				NewPassword: "new_password",
			},
			mockSetup: func(adminService *MockadminService) {
				adminService.EXPECT().UpdateAdmin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil).Times(1)
			},
			expectedCode: http.StatusOK,
		},
		{
			name:     "Incorrect password",
			username: "test",
			args: updareAdminRequest{
				Password:    "password",
				NewPassword: "new_password",
			},
			mockSetup: func(adminService *MockadminService) {
				adminService.EXPECT().UpdateAdmin(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(admin.ErrIDUsed).Times(1)
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

			req := httptest.NewRequest(http.MethodPost, "/admin/"+tt.username, bytes.NewBuffer(reqBody))
			req = mux.SetURLVars(req, map[string]string{
				AdminUsernameParam: tt.username,
			})
			res := httptest.NewRecorder()

			handler := NewHandler(mockService)
			handler.UpdateAdmin(t.Context(), res, req)

			assert.Equal(t, tt.expectedCode, res.Code)
		})
	}
}
