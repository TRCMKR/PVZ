package models

import (
	"fmt"
	"strings"
	"time"

	"github.com/Rhymond/go-money"
)

// StatusType is a type for orders in db
type StatusType uint

const (
	// StoredOrder is a status for stored orders
	StoredOrder StatusType = iota + 1

	// GivenOrder is a status for given orders
	GivenOrder

	// ReturnedOrder is a status for returned orders
	ReturnedOrder

	// DeletedOrder is a status for deleted orders
	DeletedOrder
)

// Order represents an order in the system
// @Description Order structure represents an order in the system
type Order struct {
	// @Description Unique ID of the order
	// @Example 123
	ID int `db:"id" json:"id"`

	// @Description ID of the user who created the order
	// @Example 456
	UserID int `db:"user_id" json:"user_id"`

	// @Description Weight of the order in kilograms
	// @Example 5.5
	Weight float64 `db:"weight" json:"weight"`

	// @Description Price of the order
	// @Example {"amount": 100, "currency": "USD"}
	Price money.Money `db:"price" json:"price"`

	// @Description Type of packaging for the order
	// @Example "box"
	Packaging PackagingType `db:"packaging" json:"packaging"`

	// @Description Extra packaging option for the order
	// @Example "wrap"
	ExtraPackaging PackagingType `db:"extra_packaging" json:"extra_packaging"`

	// @Description Current status of the order (e.g., 'stored', 'given', etc.)
	// @Example "stored"
	Status StatusType `db:"status" json:"status"`

	// @Description Date when the order is expected to arrive
	// @Example "2025-03-10T10:00:00Z"
	ArrivalDate time.Time `db:"arrival_date" json:"arrival_date"`

	// @Description Date when the order will expire
	// @Example "2025-03-15T10:00:00Z"
	ExpiryDate time.Time `db:"expiry_date" json:"expiry_date"`

	// @Description The last date when the order was modified
	// @Example "2025-03-09T10:00:00Z"
	LastChange time.Time `db:"last_change" json:"last_change"`
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

// NewOrder creates an instance of Order
func NewOrder(id int, userID int, weight float64, price money.Money, status StatusType,
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

// GetPackagingString makes string from packagings used for this Order
func (o *Order) GetPackagingString() string {
	var result string
	if o.Packaging == NoPackaging {
		result = "none"
	} else {
		result = GetPackagingName(o.Packaging)
	}

	if o.ExtraPackaging != NoPackaging {
		result += " in " + GetPackagingName(o.ExtraPackaging)
	}

	return result
}
