//go:generate mockgen -source=admin_handler.go -destination=../mocks/service/mock_admin_service.go -package=service
package web

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/gorilla/mux"
)

type AdminHandler struct {
	adminService adminService
}

type adminService interface {
	CreateAdmin(context.Context, models.Admin) error
	GetAdminByUsername(context.Context, string) (models.Admin, error)
	UpdateAdmin(context.Context, string, string, models.Admin) error
	DeleteAdmin(context.Context, string, string) error
	ContainsUsername(context.Context, string) (bool, error)
	ContainsID(context.Context, int) (bool, error)
}

func (h *AdminHandler) CreateAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var createRequest = struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if createRequest.ID == 0 || createRequest.Username == "" || createRequest.Password == "" {
		http.Error(w, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	admin := *models.NewAdmin(createRequest.ID, createRequest.Username, createRequest.Password)

	err = h.adminService.CreateAdmin(ctx, admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) UpdateAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	adminUsername, ok := mux.Vars(r)[adminUsernameParam]
	if !ok {
		http.Error(w, errInvalidUsername.Error(), http.StatusBadRequest)

		return
	}

	var updateRequest = struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	if updateRequest.Password == "" || updateRequest.NewPassword == "" {
		http.Error(w, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	admin := *models.NewAdmin(0, adminUsername, updateRequest.NewPassword)

	err = h.adminService.UpdateAdmin(ctx, adminUsername, updateRequest.Password, admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *AdminHandler) DeleteAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	adminUsername, ok := mux.Vars(r)[adminUsernameParam]
	if !ok {
		http.Error(w, errInvalidUsername.Error(), http.StatusBadRequest)

		return
	}
	var deleteRequest = struct {
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&deleteRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	if deleteRequest.Password == "" {
		http.Error(w, errFieldsMissing.Error(), http.StatusBadRequest)
	}

	err = h.adminService.DeleteAdmin(ctx, deleteRequest.Password, adminUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
