package repository

import (
	"database/sql"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"

	"github.com/Rhymond/go-money"
)

type order struct {
	ID             int                  `db:"id"`
	UserID         int                  `db:"user_id"`
	Weight         float64              `db:"weight"`
	Price          int64                `db:"price"`
	Packaging      models.PackagingType `db:"packaging"`
	ExtraPackaging models.PackagingType `db:"extra_packaging"`
	Status         sql.NullString       `db:"status"`
	ArrivalDate    sql.NullTime         `db:"arrival_date"`
	ExpiryDate     sql.NullTime         `db:"expiry_date"`
	LastChange     sql.NullTime         `db:"last_change"`
}

func convertToRepo(someOrder *models.Order) *order {
	orderRepo := &order{
		ID:             someOrder.ID,
		UserID:         someOrder.UserID,
		Weight:         someOrder.Weight,
		Price:          someOrder.Price.Amount(),
		Packaging:      someOrder.Packaging,
		ExtraPackaging: someOrder.ExtraPackaging,
		Status:         sql.NullString{String: someOrder.Status, Valid: true},
		ArrivalDate:    sql.NullTime{Time: someOrder.ArrivalDate, Valid: true},
		ExpiryDate:     sql.NullTime{Time: someOrder.ExpiryDate, Valid: true},
		LastChange:     sql.NullTime{Time: someOrder.LastChange, Valid: true},
	}

	return orderRepo
}

func convertToModel(someOrder *order) *models.Order {
	orderModel := &models.Order{
		ID:             someOrder.ID,
		UserID:         someOrder.UserID,
		Weight:         someOrder.Weight,
		Price:          *money.New(someOrder.Price, money.RUB),
		Packaging:      someOrder.Packaging,
		ExtraPackaging: someOrder.ExtraPackaging,
	}
	if someOrder.Status.Valid {
		orderModel.Status = someOrder.Status.String
	}
	if someOrder.ArrivalDate.Valid {
		orderModel.ArrivalDate = someOrder.ArrivalDate.Time
	}
	if someOrder.ExpiryDate.Valid {
		orderModel.ExpiryDate = someOrder.ExpiryDate.Time
	}
	if someOrder.LastChange.Valid {
		orderModel.LastChange = someOrder.LastChange.Time
	}

	return orderModel
}
