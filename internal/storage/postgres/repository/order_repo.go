package repository

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v4"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
)

// OrdersRepo is a structure for orders repo
type OrdersRepo struct {
	db database
}

// NewOrdersRepo creates an instance of orders repo
func NewOrdersRepo(db database) *OrdersRepo {
	return &OrdersRepo{
		db: db,
	}
}

func selectTx(ctx context.Context, tx pgx.Tx, orders []order, selectQuery string, args ...interface{}) error {
	rows, err := tx.Query(ctx, selectQuery, args...)
	if err != nil {
		return errGetOrdersFailed
	}

	defer rows.Close()

	for rows.Next() {
		var tmp order
		err = rows.Scan(
			&tmp.ID,
			&tmp.UserID,
			&tmp.Weight,
			&tmp.Price,
			&tmp.Packaging,
			&tmp.ExtraPackaging,
			&tmp.Status,
			&tmp.ArrivalDate,
			&tmp.ExpiryDate,
			&tmp.LastChange)
		if err != nil {
			return errGetOrdersFailed
		}

		orders = append(orders, tmp)
	}

	if err = rows.Err(); err != nil {
		return errGetOrdersFailed
	}

	return nil
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
	errFindingOrder      = errors.New("failed to find order")
)

// AddOrder adds order
func (r *OrdersRepo) AddOrder(ctx context.Context, tx pgx.Tx, order models.Order) error {
	tmp := convertToRepo(&order)

	exec := r.db.Exec
	if tx != nil {
		exec = tx.Exec
	}

	_, err := exec(ctx, `
						INSERT INTO orders(id,
										   user_id,
										   weight,
										   price,
										   packaging,
										   extra_packaging,
										   status,
										   arrival_date,
										   expiry_date,
										   last_change)
						VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);
						`,
		tmp.ID, tmp.UserID, tmp.Weight, tmp.Price, tmp.Packaging, tmp.ExtraPackaging,
		tmp.Status, tmp.ArrivalDate.Time, tmp.ExpiryDate.Time, tmp.LastChange.Time)
	if err != nil {
		log.Printf("Failed to insert order %v: %v", tmp.ID, err)

		return errAddOrderFailed
	}

	return nil
}

// RemoveOrder deletes order
func (r *OrdersRepo) RemoveOrder(ctx context.Context, tx pgx.Tx, id int) error {
	someOrder, err := r.GetByID(ctx, tx, id)
	if err != nil {
		return err
	}

	if someOrder.Status == models.DeletedOrder {
		return errNoSuchOrder
	}

	exec := r.db.Exec
	if tx != nil {
		exec = tx.Exec
	}

	_, err = exec(ctx, `
						UPDATE orders 
						SET status = $1 
						WHERE id = $2
						AND status <> $3;
						`, 4, id, models.DeletedOrder)
	if err != nil {
		log.Printf("Failed to remove order %v: %v", id, errRemoveOrderFailed)

		return errRemoveOrderFailed
	}

	return nil
}

// UpdateOrder updates order
func (r *OrdersRepo) UpdateOrder(ctx context.Context, tx pgx.Tx, id int, order models.Order) error {
	exec := r.db.Exec
	if tx != nil {
		exec = tx.Exec
	}

	_, err := exec(ctx, `
						UPDATE orders
						SET user_id         = $1,
							weight          = $2,
							price           = $3,
							packaging       = $4,
							extra_packaging = $5,
							status          = $6,
							arrival_date    = $7,
							expiry_date     = $8,
							last_change     = $9
						WHERE id = $10
						`,
		order.UserID, order.Weight, order.Price.Amount(), order.Packaging, order.ExtraPackaging,
		order.Status, order.ArrivalDate, order.ExpiryDate, order.LastChange, id)
	if err != nil {
		log.Printf("Failed to update order %v: %v", id, errUpdateOrderFailed)

		return err
	}

	return nil
}

// GetByID gets order by id
func (r *OrdersRepo) GetByID(ctx context.Context, tx pgx.Tx, id int) (models.Order, error) {
	execQueryRow := r.db.ExecQueryRow
	if tx != nil {
		execQueryRow = tx.QueryRow
	}

	var someOrder order
	err := execQueryRow(ctx, `
							SELECT * 
							FROM orders 
							WHERE id = $1
							AND status <> 4
							`, id).Scan(
		&someOrder.ID,
		&someOrder.UserID,
		&someOrder.Weight,
		&someOrder.Price,
		&someOrder.Packaging,
		&someOrder.ExtraPackaging,
		&someOrder.Status,
		&someOrder.ArrivalDate,
		&someOrder.ExpiryDate,
		&someOrder.LastChange)
	if err != nil {
		log.Printf("Failed to get order %v: %v", id, errGetOrderByID)

		return models.Order{}, errGetOrderByID
	}

	return *convertToModel(&someOrder), nil
}

// GetByUserID gets orders by user id
func (r *OrdersRepo) GetByUserID(ctx context.Context, tx pgx.Tx, id int, count int) ([]models.Order, error) {
	selectFunc := r.db.Select
	if tx != nil {
		selectFunc = func(ctx context.Context, dest interface{}, selectQuery string, args ...interface{}) error {
			return selectTx(ctx, tx, dest.([]order), selectQuery, args...)
		}
	}

	var tmp []order
	var err error
	if count == 0 {
		err = selectFunc(ctx, &tmp, `
									SELECT * 
									FROM orders 
									WHERE user_id = $1 
									AND status <> 4
									ORDER BY last_change DESC
									`, id)
	} else {
		err = selectFunc(ctx, &tmp, `
									SELECT * 
									FROM orders 
									WHERE user_id = $1 
									AND status <> 4
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

// GetReturns gets all returned orders
func (r *OrdersRepo) GetReturns(ctx context.Context, tx pgx.Tx) ([]models.Order, error) {
	selectFunc := r.db.Select
	if tx != nil {
		selectFunc = func(ctx context.Context, dest interface{}, selectQuery string, args ...interface{}) error {
			return selectTx(ctx, tx, dest.([]order), selectQuery, args...)
		}
	}

	var tmp []order
	err := selectFunc(ctx, &tmp, `
								SELECT * 
								FROM orders 
								WHERE status = 3
								ORDER BY last_change DESC
								`)
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

// GetOrders gets all orders that satisfy conditions
func (r *OrdersRepo) GetOrders(ctx context.Context, tx pgx.Tx, params []query.Cond,
	count int, page int) ([]models.Order, error) {
	var tmp []order

	params = append(params, query.Cond{
		Operator: query.NotEquals,
		Field:    "status",
		Value:    4,
	})

	selectQuery, args := query.BuildSelectQuery("orders",
		query.Where(params...),
		query.OrderBy("last_change"),
		query.Desc(true),
		query.Limit(count),
		query.Offset(page*count),
	)

	selectFunc := r.db.Select
	if tx != nil {
		selectFunc = func(ctx context.Context, dest interface{}, selectQuery string, args ...interface{}) error {
			return selectTx(ctx, tx, dest.([]order), selectQuery, args...)
		}
	}

	err := selectFunc(ctx, &tmp, selectQuery, args...)

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

// OffsetGetOrders gets orders that satisfy conditions with offset
func (r *OrdersRepo) OffsetGetOrders(ctx context.Context, tx pgx.Tx, params []query.Cond,
	count int, page int, offset int) ([]models.Order, error) {
	params = append(params, query.Cond{
		Operator: query.NotEquals,
		Field:    "status",
		Value:    4,
	})

	selectQuery, args := query.BuildSelectQuery("orders",
		query.Where(params...),
		query.OrderBy("last_change"),
		query.Desc(true),
		query.Limit(count),
		query.Offset(offset+page*count),
	)

	selectFunc := r.db.Select
	if tx != nil {
		selectFunc = func(ctx context.Context, dest interface{}, selectQuery string, args ...interface{}) error {
			return selectTx(ctx, tx, dest.([]order), selectQuery, args...)
		}
	}

	var tmp []order
	err := selectFunc(ctx, &tmp, selectQuery, args...)
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

// Contains checks if order is present
func (r *OrdersRepo) Contains(ctx context.Context, tx pgx.Tx, id int) (bool, error) {
	execQueryRow := r.db.ExecQueryRow
	if tx != nil {
		execQueryRow = tx.QueryRow
	}

	var exists bool
	err := execQueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1)`, id).Scan(&exists)
	if err != nil {
		log.Printf("Failed to find order %v: %v", id, errFindingOrder)

		return false, errFindingOrder
	}

	return exists, nil
}
