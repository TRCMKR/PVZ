package admin

import (
	"context"
	"encoding/json"
	"net/http"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

// createAdminRequest represents the request body for creating an admin
// @Description Request to create a new admin user
// @Accept json
// @Produce json
// @Param request body createAdminRequest true "Create Admin Request"
// @Success 200 {string} string "Admin created successfully"
// @Failure 400 {string} string "Invalid request or missing fields"
// @Failure 500 {string} string "Internal server error"
// @Router /admin/create [post]
type createAdminRequest struct {
	ID       int    `json:"id"`       // ID is the unique identifier for the admin
	Username string `json:"username"` // Username is the name the admin will use to log in
	Password string `json:"password"` // Password is the admin's password
}

// CreateAdmin creates admin
// @Summary Create admin
// @Description Creates a new admin user
// @Tags admins
// @Accept json
// @Produce json
// @Param admin body createAdminRequest true "Admin details"
// @Success 200 {string} string "Admin created successfully"
// @Failure 400 {string} string "Invalid request or missing fields"
// @Failure 500 {string} string "Internal server error"
// @Router /admins [post]
func (h *Handler) CreateAdmin(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var createRequest = createAdminRequest{}
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
	_, _ = w.Write([]byte("admin created successfully"))
}
