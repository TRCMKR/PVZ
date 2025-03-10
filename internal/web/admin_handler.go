package web

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/gorilla/mux"
)

func (s *server) CreateAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	type createAdminRequest struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var createRequest createAdminRequest
	err := json.NewDecoder(r.Body).Decode(&createRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	if createRequest.ID == 0 || createRequest.Username == "" || createRequest.Password == "" {
		http.Error(w, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	admin := *models.NewAdmin(createRequest.ID, createRequest.Username, createRequest.Password)

	err = s.adminService.CreateAdmin(ctx, admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) UpdateAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	type updateAdminRequest struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}

	adminUsername, ok := mux.Vars(r)[adminUsernameParam]
	if !ok {
		http.Error(w, errInvalidUsername.Error(), http.StatusBadRequest)

		return
	}

	var updateRequest updateAdminRequest
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

	err = s.adminService.UpdateAdmin(ctx, adminUsername, updateRequest.Password, admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *server) DeleteAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	type deleteAdminRequest struct {
		Password string `json:"password"`
	}

	adminUsername, ok := mux.Vars(r)[adminUsernameParam]
	if !ok {
		http.Error(w, errInvalidUsername.Error(), http.StatusBadRequest)

		return
	}
	var deleteRequest deleteAdminRequest
	err := json.NewDecoder(r.Body).Decode(&deleteRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}
	if deleteRequest.Password == "" {
		http.Error(w, errFieldsMissing.Error(), http.StatusBadRequest)
	}

	err = s.adminService.DeleteAdmin(ctx, deleteRequest.Password, adminUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
