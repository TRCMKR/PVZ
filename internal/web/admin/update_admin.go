package admin

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/gorilla/mux"
)

func (h *Handler) UpdateAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	adminUsername, ok := mux.Vars(r)[AdminUsernameParam]
	if !ok {
		http.Error(w, ErrNoUsername.Error(), http.StatusBadRequest)

		return
	}

	var updateRequest = struct {
		Password    string `json:"password"`
		NewPassword string `json:"new_password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&updateRequest)
	if err != nil {
		http.Error(w, ErrFieldsMissing.Error(), http.StatusBadRequest)

		return
	}
	if updateRequest.Password == "" || updateRequest.NewPassword == "" {
		http.Error(w, ErrFieldsMissing.Error(), http.StatusBadRequest)

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
