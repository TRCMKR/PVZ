package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
)

type Order struct {
	ID             int           `db:"id"`
	UserID         int           `db:"user_id"`
	Weight         float64       `db:"weight"`
	Price          money.Money   `db:"price"`
	Packaging      PackagingType `db:"packaging"`
	ExtraPackaging PackagingType `db:"extra_packaging"`
	Status         string        `db:"status"`
	ArrivalDate    time.Time     `db:"arrival_date"`
	ExpiryDate     time.Time     `db:"expiry_date"`
	LastChange     time.Time     `db:"last_change"`
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

const (
	dateLayout = "2006.01.02 15:04:05"
)

func NewOrder(id int, userID int, weight float64, price money.Money, status string,
	arrivalDate time.Time, expiryDate time.Time, lastChange time.Time) *Order {
	return &Order{
		ID:             id,
		UserID:         userID,
		Weight:         weight,
		Price:          price,
		Packaging:      NoPackaging,
		ExtraPackaging: NoPackaging,
		Status:         status,
		ArrivalDate:    arrivalDate,
		ExpiryDate:     expiryDate,
		LastChange:     lastChange,
	}
}

func (o *Order) String() string {
	sb := strings.Builder{}
	rowFormat := fmt.Sprintf("OID%%%dd | UID%%%dd | WGHT%%%d.2f | PRC%%%ds | PKG%%%ds | STAT%%%ds | LCHAN%%%ds",
		orderIDWidth, userIDWidth, weightWidth, priceWidth, packagingWidth, statusWidth, dateWidth)

	_, err := fmt.Fprintf(&sb, rowFormat, o.ID, o.UserID,
		o.Weight, o.Price.Display(), o.GetPackagingString(), o.Status, o.LastChange.Format(dateLayout))
	if err != nil {
		sb.WriteString(err.Error())
	}

	return sb.String()
}

func (o *Order) GetPackagingString() string {
	var result string
	if o.Packaging == NoPackaging {
		result = "none"
	} else {
		result = getPackagingName(o.Packaging)
	}

	if o.ExtraPackaging != NoPackaging {
		result += " in " + getPackagingName(o.ExtraPackaging)
	}

	return result
}
