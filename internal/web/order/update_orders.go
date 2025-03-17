package order

import (
	"context"
	"encoding/json"
	"net/http"
)

type processOrderRequest struct {
	UserID   int    `json:"user_id"`
	OrderIDs []int  `json:"order_ids"`
	Action   string `json:"action"`
}

func (h *Handler) UpdateOrders(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var processRequest processOrderRequest

	err := json.NewDecoder(r.Body).Decode(&processRequest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)

		return
	}

	if len(processRequest.OrderIDs) == 0 || processRequest.Action == "" {
		http.Error(w, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	var response = struct {
		Failed int `json:"failed"`
	}{}
	response.Failed, err = h.OrderService.ProcessOrders(ctx, processRequest.UserID,
		processRequest.OrderIDs, processRequest.Action)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	data, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}
