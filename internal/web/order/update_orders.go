package order

import (
	"context"
	"encoding/json"
	"net/http"
)

// processOrderRequest represents the request body for the UpdateOrder endpoint
// @Description Request to process orders by action and order IDs
// @Accept json
// @Produce json
// @Param request body processOrderRequest true "Process Orders Request"
// @Success 200 {object} struct{ Failed int } "Number of failed orders"
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /orders/update [post]
type processOrderRequest struct {
	OrderID int    `json:"id"`
	UserID  int    `json:"user_id"`
	Action  string `json:"action"`
}

// UpdateOrder updates the orders based on the provided request data
// @Security BasicAuth
// @Summary Process orders
// @Description Processes the given orders based on the action and order IDs provided
// @Tags orders
// @Accept json
// @Produce json
// @Param request body processOrderRequest true "Process Orders Request"
// @Success 200 {object} processOrderRequest
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Internal server error"
// @Router /orders/process [post]
func (h *Handler) UpdateOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var processRequest processOrderRequest

	err := json.NewDecoder(r.Body).Decode(&processRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if processRequest.OrderID == 0 || processRequest.Action == "" {
		http.Error(w, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	err = h.OrderService.ProcessOrder(ctx, processRequest.UserID,
		processRequest.OrderID, processRequest.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("success"))
}
