package order

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
}
