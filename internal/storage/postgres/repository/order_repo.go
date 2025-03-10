package repository

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
)

type OrderRepo struct {
	db database
}

func NewOrderRepo(db database) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

var (
	errAddOrderFailed    = errors.New("failed to add order")
	errRemoveOrderFailed = errors.New("failed to remove order")
	errUpdateOrderFailed = errors.New("failed to update order")
	errGetOrderByID      = errors.New("failed to get order by id")
	errGetOrderByUserID  = errors.New("failed to get order by user id")
	errGetOrdersFailed   = errors.New("failed to get orders")
	errGetReturnsFailed  = errors.New("failed to get order returns")
	errNoSuchOrder       = errors.New("no such order")
)

func (r *OrderRepo) AddOrder(ctx context.Context, order models.Order) error {
	tmp := convertToRepo(&order)
	_, err := r.db.Exec(ctx, `
INSERT INTO orders(id, user_id, weight, price, packaging, extra_packaging,status, arrival_date, expiry_date, last_change)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`,
		tmp.ID, tmp.UserID, tmp.Weight, tmp.Price, tmp.Packaging, tmp.ExtraPackaging,
		tmp.Status, tmp.ArrivalDate, tmp.ExpiryDate, tmp.LastChange)

	if err != nil {
		log.Printf("Failed to insert order %v: %v", tmp.ID, errAddOrderFailed)

		return errAddOrderFailed
	}

	return nil
}

func (r *OrderRepo) RemoveOrder(ctx context.Context, id int) error {
	someOrder, err := r.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if someOrder.Status == "deleted" {
		return errNoSuchOrder
	}

	_, err = r.db.Exec(ctx, `
UPDATE orders 
SET status = $1 
WHERE id = $2
`, "deleted", id)

	if err != nil {
		log.Printf("Failed to remove order %v: %v", id, errRemoveOrderFailed)

		return errRemoveOrderFailed
	}

	return nil
}

func (r *OrderRepo) UpdateOrder(ctx context.Context, id int, order models.Order) error {
	_, err := r.db.Exec(ctx, `
UPDATE orders 
SET user_id = $1, weight = $2, price = $3, packaging = $4, extra_packaging = $5, status = $6, 
    arrival_date = $7, expiry_date = $8, last_change = $9 
WHERE id = $10
`,
		order.UserID, order.Weight, order.Price.Amount(), order.Packaging, order.ExtraPackaging,
		order.Status, order.ArrivalDate, order.ExpiryDate, order.LastChange, id)

	if err != nil {
		log.Printf("Failed to update order %v: %v", id, errUpdateOrderFailed)

		return errUpdateOrderFailed
	}

	return nil
}

func (r *OrderRepo) GetByID(ctx context.Context, id int) (models.Order, error) {
	var someOrder order
	err := r.db.Get(ctx, &someOrder, `
SELECT * 
FROM orders 
WHERE id = $1
`, id)

	if err != nil {
		log.Printf("Failed to get order %v: %v", id, errGetOrderByID)

		return models.Order{}, errGetOrderByID
	}

	return *convertToModel(&someOrder), nil
}

func (r *OrderRepo) GetByUserID(ctx context.Context, id int, count int) ([]models.Order, error) {
	var tmp []order
	var err error
	if count == 0 {
		err = r.db.Select(ctx, &tmp, `
SELECT * 
FROM orders 
WHERE user_id = $1 
ORDER BY last_change DESC
`, id)
	} else {
		err = r.db.Select(ctx, &tmp, `
SELECT * 
FROM orders 
WHERE user_id = $1 
ORDER BY last_change DESC
LIMIT $2
`, id, count)
	}

	if err != nil {
		log.Printf("Failed to get order by user %v: %v", id, errGetOrderByUserID)

		return nil, errGetOrderByUserID
	}

	orders := make([]models.Order, 0, len(tmp))
	for x := range tmp {
		orders = append(orders, *convertToModel(&tmp[x]))
	}

	return orders, nil
}

func (r *OrderRepo) GetReturns(ctx context.Context) ([]models.Order, error) {
	var tmp []order
	err := r.db.Select(ctx, &tmp, "SELECT * FROM orders WHERE status = 'returned' ORDER BY last_change DESC")

	if err != nil {
		log.Printf("Failed to get order by user id: %v", errGetReturnsFailed)

		return nil, errGetReturnsFailed
	}

	orders := make([]models.Order, 0, len(tmp))
	for x := range tmp {
		orders = append(orders, *convertToModel(&tmp[x]))
	}

	return orders, nil
}

func (r *OrderRepo) makeCondition(k string, paramInd int) string {
	switch k {
	case "arrival_date_from":
		return fmt.Sprintf("arrival_date >= $%d", paramInd)
	case "arrival_date_to":
		return fmt.Sprintf("arrival_date <= $%d", paramInd)
	case "expiry_date_from":
		return fmt.Sprintf("expiry_date >= $%d", paramInd)
	case "expiry_date_to":
		return fmt.Sprintf("expiry_date <= $%d", paramInd)
	case "weight_from":
		return fmt.Sprintf("weight >= $%d", paramInd)
	case "weight_to":
		return fmt.Sprintf("weight <= $%d", paramInd)
	case "price_from":
		return fmt.Sprintf("price >= $%d", paramInd)
	case "price_to":
		return fmt.Sprintf("price <= $%d", paramInd)
	default:
		return fmt.Sprintf("%s = $%d", k, paramInd)
	}
}

func (r *OrderRepo) makePagination(args []interface{}, paramInd *int, count int, page int) ([]interface{}, string) {
	var sb strings.Builder
	if count > 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT $%d", *paramInd))
		args = append(args, count)
		*paramInd++
		if page > 0 {
			sb.WriteString(fmt.Sprintf(" OFFSET $%d", *paramInd))
			args = append(args, page*count)
			*paramInd++
		}
	}

	return args, sb.String()
}

func (r *OrderRepo) GetOrders(ctx context.Context, params map[string]string,
	count int, page int) ([]models.Order, error) {
	var tmp []order
	var sb strings.Builder

	sb.WriteString("SELECT * FROM orders")
	args := make([]interface{}, 0, len(params))
	paramInd := 1
	sb.WriteString(" WHERE status != 'deleted'")

	for param, value := range params {
		if value == "" {
			continue
		}
		sb.WriteString(" AND ")
		sb.WriteString(r.makeCondition(param, paramInd))
		args = append(args, value)
		paramInd++
	}

	sb.WriteString(" ORDER BY last_change DESC")

	var pagination string
	args, pagination = r.makePagination(args, &paramInd, count, page)
	sb.WriteString(pagination)
	sb.WriteByte(';')

	err := r.db.Select(ctx, &tmp, sb.String(), args...)

	if err != nil {
		log.Printf("Failed to get orders %v: %v", tmp, errGetOrdersFailed)

		return nil, errGetOrdersFailed
	}

	orders := make([]models.Order, 0, len(tmp))
	for x := range tmp {
		orders = append(orders, *convertToModel(&tmp[x]))
	}

	return orders, nil
}

func (r *OrderRepo) Contains(ctx context.Context, id int) bool {
	var exists bool
	_ = r.db.Get(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)`, id)

	return exists
}
