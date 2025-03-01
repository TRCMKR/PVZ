package order

import (
	"fmt"
	"strings"

	"github.com/Rhymond/go-money"
)

type Order struct {
	ID             int
	UserID         int
	Weight         float64
	Price          money.Money
	Packaging      packagingType
	ExtraPackaging packagingType
	Status         string
	ArrivalDate    string
	ExpiryDate     string
	LastChange     string
}

const (
	orderIDWidth   = 12
	userIDWidth    = 10
	weightWidth    = 8
	priceWidth     = 12
	packagingWidth = 13
	statusWidth    = 10
	dateWidth      = 22
)

func NewOrder(id int, userID int, weight float64, price money.Money, status string,
	arrivalDate string, expiryDate string, lastChange string) *Order {
	return &Order{
		ID:             id,
		UserID:         userID,
		Weight:         weight,
		Price:          price,
		Packaging:      noPackaging,
		ExtraPackaging: noPackaging,
		ArrivalDate:    arrivalDate,
		Status:         status,
		ExpiryDate:     expiryDate,
		LastChange:     lastChange,
	}
}

func (o *Order) String() string {
	sb := strings.Builder{}
	rowFormat := fmt.Sprintf("OID%%%dd | UID%%%dd | WGHT%%%d.2f | PRC%%%ds | PKG%%%ds | STAT%%%ds | LCHAN%%%ds",
		orderIDWidth, userIDWidth, weightWidth, priceWidth, packagingWidth, statusWidth, dateWidth)

	packgingString := formPackagingString(o.Packaging, o.ExtraPackaging)

	_, err := fmt.Fprintf(&sb, rowFormat, o.ID, o.UserID,
		o.Weight, o.Price.Display(), packgingString, o.Status, o.LastChange)
	if err != nil {
		sb.WriteString(err.Error())
	}

	return sb.String()
}
