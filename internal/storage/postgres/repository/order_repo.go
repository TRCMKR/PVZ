package repository

import (
	"context"
	"fmt"
	"strings"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/storage/postgres"
)

type OrderRepo struct {
	db postgres.Database
}

func NewOrderRepo(db postgres.Database) *OrderRepo {
	return &OrderRepo{
		db: db,
	}
}

func (r *OrderRepo) AddOrder(ctx context.Context, someOrder models.Order) {
	tmp := convertToRepo(&someOrder)
	_, _ = r.db.Exec(ctx, `INSERT INTO orders(id, user_id, weight, price, packaging, extra_packaging,
		status, arrival_date, expiry_date, last_change) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		tmp.ID, tmp.UserID, tmp.Weight, tmp.Price, tmp.Packaging, tmp.ExtraPackaging,
		tmp.Status, tmp.ArrivalDate, tmp.ExpiryDate, tmp.LastChange)
}

func (r *OrderRepo) RemoveOrder(ctx context.Context, id int) {
	_, _ = r.db.Exec(ctx, "UPDATE orders SET status = $1 WHERE id = $2", "deleted", id)
}

func (r *OrderRepo) UpdateOrder(ctx context.Context, id int, order models.Order) {
	_, _ = r.db.Exec(ctx, `UPDATE orders 
		SET user_id = $1, weight = $2, price = $3, packaging = $4, extra_packaging = $5, 
		    status = $6, arrival_date = $7, expiry_date = $8, last_change = $9 
		WHERE id = $10`,
		order.UserID, order.Weight, order.Price.Amount(), order.Packaging, order.ExtraPackaging,
		order.Status, order.ArrivalDate, order.ExpiryDate, order.LastChange, id)
}

func (r *OrderRepo) GetByID(ctx context.Context, id int) models.Order {
	var someOrder order
	_ = r.db.Get(ctx, &someOrder, "SELECT * FROM orders WHERE id = $1", id)

	return *convertToModel(&someOrder)
}

func (r *OrderRepo) GetByUserID(ctx context.Context, id int, count int) []models.Order {
	var tmp []order
	if count == 0 {
		_ = r.db.Select(ctx, &tmp, "SELECT * FROM orders WHERE user_id = $1 ORDER BY last_change DESC", id)
	} else {
		_ = r.db.Select(ctx, &tmp, "SELECT * FROM orders WHERE user_id = $1 ORDER BY last_change DESC LIMIT $2", id, count)
	}
	orders := make([]models.Order, 0, len(tmp))
	for x := range tmp {
		orders = append(orders, *convertToModel(&tmp[x]))
	}

	return orders
}

func (r *OrderRepo) GetReturns(ctx context.Context) []models.Order {
	var tmp []order
	_ = r.db.Select(ctx, &tmp, "SELECT * FROM orders WHERE status = 'returned' ORDER BY last_change DESC")
	orders := make([]models.Order, 0, len(tmp))
	for x := range tmp {
		orders = append(orders, *convertToModel(&tmp[x]))
	}

	return orders
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

func (r *OrderRepo) GetOrders(ctx context.Context, params map[string]string, count int, page int) []models.Order {
	var tmp []order
	var sb strings.Builder

	sb.WriteString("SELECT * FROM orders")
	args := make([]interface{}, 0, len(params))
	paramInd := 1

	for k, v := range params {
		if v == "" {
			continue
		}
		if len(args) == 0 {
			sb.WriteString(" WHERE ")
		} else {
			sb.WriteString(" AND ")
		}
		sb.WriteString(r.makeCondition(k, paramInd))
		args = append(args, v)
		paramInd++
	}

	if paramInd == 1 {
		sb.WriteString(" WHERE status != 'deleted' ")
	} else {
		sb.WriteString(" AND status != 'deleted' ")
	}

	sb.WriteString(" ORDER BY last_change DESC")

	var pagination string
	args, pagination = r.makePagination(args, &paramInd, count, page)
	sb.WriteString(pagination)
	sb.WriteByte(';')

	_ = r.db.Select(ctx, &tmp, sb.String(), args...)
	orders := make([]models.Order, 0, len(tmp))
	for x := range tmp {
		orders = append(orders, *convertToModel(&tmp[x]))
	}

	return orders
}

func (r *OrderRepo) Save(ctx context.Context) error {
	return nil
}

func (r *OrderRepo) Contains(ctx context.Context, id int) bool {
	var exists bool
	_ = r.db.Get(ctx, &exists, "SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)", id)

	return exists
}
