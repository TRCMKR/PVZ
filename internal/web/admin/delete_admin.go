package admin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func (h *Handler) DeleteAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	adminUsername, ok := mux.Vars(r)[AdminUsernameParam]
	if !ok {
		http.Error(w, ErrNoUsername.Error(), http.StatusBadRequest)

		return
	}
	var deleteRequest = struct {
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(r.Body).Decode(&deleteRequest)
	if err != nil {
		http.Error(w, ErrFieldsMissing.Error(), http.StatusBadRequest)

		return
	}
	if deleteRequest.Password == "" {
		http.Error(w, ErrFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	err = h.adminService.DeleteAdmin(ctx, deleteRequest.Password, adminUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
