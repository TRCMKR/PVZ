package models

import "homework/order"

type Storage interface {
	AcceptOrder(orderID string, userID string, expiryDate string) error
	AcceptOrders(path string) (int, error)
	ReturnOrder(orderID string) error
	ProcessOrders(userID string, orderIDs []string, action string) (int, error)
	UserOrders(args ...string) ([]*order.Order, error)
	Returns() []*order.Order
	OrderHistory() []*order.Order
	Save() error
}
