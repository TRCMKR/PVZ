package order

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/Rhymond/go-money"
)

// CreateOrder handles the creation of an order
// @Security BasicAuth
// @Summary Create a new order
// @Description Creates a new order based on the provided order details and validates the fields
// @Tags orders
// @Accept  json
// @Produce  json
// @Param order body createOrderRequest true "Order details"
// @Success 200 {string} string "Success"
// @Failure 400 {string} string "Invalid JSON format"
// @Failure 400 {string} string "Missing required fields"
// @Failure 400 {string} string "Invalid packaging"
// @Router /orders [post]
func (h *Handler) CreateOrder(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var order = createOrderRequest{}

	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, errWrongJsonFormat.Error(), http.StatusBadRequest)

		return
	}

	if order.ID == 0 || order.UserID == 0 || order.Weight == 0 || order.Price.IsZero() || order.ExpiryDate.IsZero() {
		http.Error(w, errFieldsMissing.Error(), http.StatusBadRequest)

		return
	}

	packagings := make([]models.Packaging, 0, 2)
	packaging, err := getPackaging(models.GetPackagingName(models.PackagingType(order.Packaging)))
	extraPackaging, errExtra := getPackaging(models.GetPackagingName(models.PackagingType(order.ExtraPackaging)))
	if err != nil || errExtra != nil {
		http.Error(w, errNoSuchPackaging.Error(), http.StatusBadRequest)

		return
	}

	packagings = append(packagings, packaging)
	packagings = append(packagings, extraPackaging)

	err = h.OrderService.AcceptOrder(ctx, order.ID, order.UserID, order.Weight, order.Price, order.ExpiryDate, packagings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("success"))
}

func getPackaging(packagingStr string) (models.Packaging, error) {
	var packaging models.Packaging
	if packagingStr != "" {
		packaging = models.GetPackaging(packagingStr)
		if packaging == nil {
			return nil, errNoSuchPackaging
		}
	} else {
		return nil, errNoSuchPackaging
	}

	return packaging, nil
}

type createOrderRequest struct {
	ID             int         `json:"id"`
	UserID         int         `json:"user_id"`
	Weight         float64     `json:"weight"`
	Price          money.Money `json:"price"`
	Packaging      int         `json:"packaging"`
	ExtraPackaging int         `json:"extra_packaging"`
	Status         int         `json:"status"`
	ExpiryDate     time.Time   `json:"expiry_date"`
}
