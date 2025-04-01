package models

import (
	"fmt"
	"time"
)

// Log ...
type Log struct {
	ID      int       `json:"id"       db:"id"`
	OrderID int       `json:"order_id" db:"order_id"`
	AdminID int       `json:"admin_id" db:"admin_id"`
	Message string    `json:"message"  db:"message"`
	Date    time.Time `json:"date"     db:"date"`
	URL     string    `json:"url"      db:"url"`
	Method  string    `json:"method"   db:"method"`
	Status  int       `json:"status"   db:"status"`
}

// NewLog ...
func NewLog(orderID int, adminID int, message string, url string, method string, status int) *Log {
	return &Log{
		ID:      0,
		OrderID: orderID,
		AdminID: adminID,
		Message: message,
		Date:    time.Now(),
		URL:     url,
		Method:  method,
		Status:  status,
	}
}

func (l *Log) String() string {
	return fmt.Sprintf("%s\nOrder %d, admin %d:\nResponse: %s\nPath: %s\nMethod: %s\nStatus: %d\n",
		l.Date, l.OrderID, l.AdminID, l.Message, l.URL, l.Method, l.Status)
}
