package order

import (
	"fmt"
	"strings"
)

type Order struct {
	ID          int
	UserID      int
	Status      string
	ArrivalDate string
	ExpiryDate  string
	LastChange  string
}

const (
	orderIDWidth = 15
	userIDWidth  = 10
	statusWidth  = 10
	dateWidth    = 22
)

func (o *Order) String() string {
	sb := strings.Builder{}
	rowFormat := fmt.Sprintf("OID%%%dd | UID%%%dd | STAT%%%ds | LCHAN%%%ds", orderIDWidth, userIDWidth, statusWidth, dateWidth)
	_, err := fmt.Fprintf(&sb, rowFormat, o.ID, o.UserID, o.Status, o.LastChange)
	if err != nil {
		sb.WriteString(err.Error())
	}

	return sb.String()
}
