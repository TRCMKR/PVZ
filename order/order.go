package order

import (
	"fmt"
	"strings"
)

type Order struct {
	OrderID     string
	UserId      string
	Status      string
	ArrivalDate string
	ExpiryDate  string
	LastChange  string
}

const (
	orderIDWidth = 15
	userIDWidth  = 10
	statusWidth  = 10
	arrivalWidth = 22
)

func (o *Order) String() string {
	sb := strings.Builder{}
	rowFormat := fmt.Sprintf("%%%ds %%%ds %%%ds %%%ds", orderIDWidth, userIDWidth, statusWidth, arrivalWidth)
	fmt.Fprintf(&sb, rowFormat, o.OrderID, o.UserId, o.Status, o.LastChange)
	//sb.WriteString(o.OrderID)
	//sb.WriteByte(' ')
	//sb.WriteString(o.UserId)
	//sb.WriteByte(' ')
	//sb.WriteString(o.Status)
	//sb.WriteByte(' ')
	//sb.WriteString(o.LastChange)

	return sb.String()
}
