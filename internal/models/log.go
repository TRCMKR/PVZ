package models

import (
	"fmt"
	"time"
)

// Log is a structure that contains all log data and necessary information to make a job from it
type Log struct {
	ID           int       `db:"id" json:"id"`
	OrderID      int       `db:"order_id" json:"order_id"`
	AdminID      int       `db:"admin_id" json:"admin_id"`
	Message      string    `db:"message" json:"message"`
	Date         time.Time `db:"date" json:"date"`
	URL          string    `db:"url" json:"url"`
	Method       string    `db:"method" json:"method"`
	Status       int       `db:"status" json:"status"`
	JobStatus    int       `db:"job_status" json:"job_status"`
	AttemptsLeft int       `db:"attempts_left" json:"attempts_left"`
	UpdatedAt    time.Time `db:"updated_at" json:"updated_at"`
}

// NewLog creates an instance of Log
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
