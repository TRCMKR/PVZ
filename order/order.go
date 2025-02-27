package order

import (
	"fmt"
	"strings"
)

type Order struct {
	ID          int
	UserID      int
	Weight      float64
	Price       float64
	Packaging   string
	Status      string
	ArrivalDate string
	ExpiryDate  string
	LastChange  string
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

func (o *Order) String() string {
	sb := strings.Builder{}
	rowFormat := fmt.Sprintf("OID%%%dd | UID%%%dd | WGHT%%%d.2f | PRC%%%d.2f | PKG%%%ds | STAT%%%ds | LCHAN%%%ds",
		orderIDWidth, userIDWidth, weightWidth, priceWidth, packagingWidth, statusWidth, dateWidth)

	_, err := fmt.Fprintf(&sb, rowFormat, o.ID, o.UserID, o.Weight, o.Price, o.Packaging, o.Status, o.LastChange)
	if err != nil {
		sb.WriteString(err.Error())
	}

	return sb.String()
}
