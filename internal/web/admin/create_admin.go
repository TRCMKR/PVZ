package admin

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

func (h *Handler) CreateAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
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
		http.Error(w, ErrFieldsMissing.Error(), http.StatusBadRequest)

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
