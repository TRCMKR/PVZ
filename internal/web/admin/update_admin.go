package admin

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/gorilla/mux"
)

// updateRequest represents the request body for updating an admin's password
// @Description Request to update the password of an admin by providing the old and new passwords
// @Accept json
// @Produce json
// @Param request body updateRequest true "Update Admin Request"
// @Success 200 {string} string "Admin password updated successfully"
// @Failure 400 {string} string "Invalid request or missing fields"
// @Failure 500 {string} string "Internal server error"
// @Router /admin/update/{username} [put]
type updateRequest struct {
	Password    string `json:"password"`     // Old password of the admin
	NewPassword string `json:"new_password"` // New password for the admin
}

// UpdateAdmin updates an admin's password
// @Summary Update admin's password
// @Description Update the password of an admin by providing the old and new passwords
// @Tags admins
// @Accept json
// @Produce json
// @Param username path string true "Admin Username" // Param for username from URL
// @Param request body updateRequest true "Update Admin Request"
// @Success 200 {string} string "Admin password updated successfully"
// @Failure 400 {string} string "Invalid request or missing fields"
// @Failure 500 {string} string "Internal server error"
// @Router /admins/{username} [post]
func (h *Handler) UpdateAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	adminUsername, ok := mux.Vars(r)[AdminUsernameParam]
	if !ok {
		http.Error(w, ErrNoUsername.Error(), http.StatusBadRequest)

		return
	}

	var request updateRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, ErrFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	if request.Password == "" || request.NewPassword == "" {
		http.Error(w, ErrFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	admin := *models.NewAdmin(0, adminUsername, request.NewPassword)

	err = h.adminService.UpdateAdmin(ctx, adminUsername, request.Password, admin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
