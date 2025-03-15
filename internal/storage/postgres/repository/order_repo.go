package repository

import (
	"context"
	"errors"
	"log"

	"gitlab.ozon.dev/alexplay1224/homework/internal/models"
	"gitlab.ozon.dev/alexplay1224/homework/internal/query"
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
	errFindingOrder      = errors.New("failed to find order")
)

func (r *OrderRepo) AddOrder(ctx context.Context, order models.Order) error {
	tmp := convertToRepo(&order)
	_, err := r.db.Exec(ctx, `
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

	if someOrder.Status == models.DeletedOrder {
		return errNoSuchOrder
	}

	_, err = r.db.Exec(ctx, `
							UPDATE orders 
							SET status = $1 
							WHERE id = $2
							AND status <> 4
							`, 4, id)

	if err != nil {
		log.Printf("Failed to remove order %v: %v", id, errRemoveOrderFailed)

		return errRemoveOrderFailed
	}

	return nil
}

func (r *OrderRepo) UpdateOrder(ctx context.Context, id int, order models.Order) error {
	_, err := r.db.Exec(ctx, `
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
									AND status <> 4
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
									AND status <> 4
									ORDER BY last_change DESC
									`, id)
	} else {
		err = r.db.Select(ctx, &tmp, `
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

func (r *OrderRepo) GetReturns(ctx context.Context) ([]models.Order, error) {
	var tmp []order
	err := r.db.Select(ctx, &tmp, `
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

func (r *OrderRepo) GetOrders(ctx context.Context, params []query.Cond,
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

	err := r.db.Select(ctx, &tmp, selectQuery, args...)

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

func (r *OrderRepo) Contains(ctx context.Context, id int) (bool, error) {
	var exists bool
	err := r.db.Get(ctx, &exists, `SELECT EXISTS(SELECT 1 FROM orders WHERE id = $1 AND status <> 4)`, id)
	if err != nil {
		log.Printf("Failed to find order %v: %v", id, errFindingOrder)

		return false, errFindingOrder
	}

	return exists, nil
}
