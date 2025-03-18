package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// DeleteOrder handles the deletion of an order
// @Security BasicAuth
// @Summary Delete an order by its ID
// @Description Deletes the specified order and returns success or error response
// @Tags orders
// @Accept  json
// @Produce  json
// @Param orderID path int true "Order ID"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Invalid Order ID"
// @Failure 500 {string} string "Internal Server Error"
// @Router /orders/{orderID} [delete]
func (h *Handler) DeleteOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.Atoi(mux.Vars(r)[OrderIDParam])
	if err != nil {
		http.Error(w, errInvalidOrderID.Error(), http.StatusBadRequest)

		return
	}

	err = h.OrderService.ReturnOrder(ctx, orderID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("success"))
}
