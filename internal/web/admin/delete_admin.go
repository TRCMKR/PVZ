package admin

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// deleteRequest represents the request body for deleting an admin
// @Description Request to delete an admin by providing the username and password
// @Accept json
// @Produce json
// @Param request body deleteRequest true "Delete Admin Request"
// @Success 200 {string} string "Admin deleted successfully"
// @Failure 400 {string} string "Invalid request or missing fields"
// @Failure 500 {string} string "Internal server error"
// @Router /admin/delete/{username} [delete]
type deleteRequest struct {
	Password string `json:"password"` // Password is required to confirm the deletion
}

// DeleteAdmin deletes an admin user
// @Summary Delete an admin
// @Description Delete an admin by providing the username and password for confirmation
// @Tags admins
// @Accept json
// @Produce json
// @Param username path string true "Admin Username" // Param for username from URL
// @Param request body deleteRequest true "Delete Admin Request"
// @Success 200 {string} string "Admin deleted successfully"
// @Failure 400 {string} string "Invalid request or missing fields"
// @Failure 500 {string} string "Internal server error"
// @Router /admins/{username} [delete]
func (h *Handler) DeleteAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	adminUsername, ok := mux.Vars(r)[AdminUsernameParam]
	if !ok {
		http.Error(w, ErrNoUsername.Error(), http.StatusBadRequest)

		return
	}
	var request deleteRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, ErrFieldsMissing.Error(), http.StatusBadRequest)

		return
	}
	if request.Password == "" {
		http.Error(w, ErrFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	err = h.adminService.DeleteAdmin(ctx, request.Password, adminUsername)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
}
