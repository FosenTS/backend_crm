package postgre

import (
	"backend_crm/internal/model"
	"backend_crm/internal/repository/orders"
	"context"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) orders.Repository {
	return &repository{db: db}
}

func (r *repository) Save(ctx context.Context, newOrder *model.NewOrder) error {
	query := `
		INSERT INTO orders (product_id, phone, email, description, status)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		newOrder.ProductId,
		newOrder.Phone,
		newOrder.Email,
		newOrder.Description,
		newOrder.Status,
	)

	return err
}

func (r *repository) GetAll(ctx context.Context) ([]*model.Order, error) {
	query := `
		SELECT o.order_id, o.phone, o.email, o.description, o.status,
			   p.product_id, p.name, p.weight, p.description
		FROM orders o
		JOIN products p ON o.product_id = p.product_id
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		var product model.Product
		err := rows.Scan(
			&order.OrderId,
			&order.Phone,
			&order.Email,
			&order.Description,
			&order.Status,
			&product.ProductId,
			&product.Name,
			&product.Weigth,
			&product.Description,
		)
		if err != nil {
			return nil, err
		}
		order.Product = product
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *repository) GetByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error) {
	return r.getOrdersByFilter(ctx, "status = $1", status)
}

func (r *repository) GetByStatusAndPhone(ctx context.Context, status model.OrderStatus, phone string) ([]*model.Order, error) {
	return r.getOrdersByFilter(ctx, "status = $1 AND phone = $2", status, phone)
}

func (r *repository) GetByStatusAndEmail(ctx context.Context, status model.OrderStatus, email string) ([]*model.Order, error) {
	return r.getOrdersByFilter(ctx, "status = $1 AND email = $2", status, email)
}

func (r *repository) GetByStatusAndPhoneAndEmail(ctx context.Context, status model.OrderStatus, phone string, email string) ([]*model.Order, error) {
	return r.getOrdersByFilter(ctx, "status = $1 AND phone = $2 AND email = $3", status, phone, email)
}

func (r *repository) GetByUserIdAndStatus(ctx context.Context, userId string, status model.OrderStatus) ([]*model.Order, error) {
	return r.getOrdersByFilter(ctx, "user_id = $1 AND status = $2", userId, status)
}

func (r *repository) GetByUserIdAndStatusAndPhone(ctx context.Context, userId string, status model.OrderStatus, phone string) ([]*model.Order, error) {
	return r.getOrdersByFilter(ctx, "user_id = $1 AND status = $2 AND phone = $3", userId, status, phone)
}

func (r *repository) GetByUserIdAndStatusAndEmail(ctx context.Context, userId string, status model.OrderStatus, email string) ([]*model.Order, error) {
	return r.getOrdersByFilter(ctx, "user_id = $1 AND status = $2 AND email = $3", userId, status, email)
}

func (r *repository) GetByUserIdAndStatusAndPhoneAndEmail(ctx context.Context, userId string, status model.OrderStatus, phone string, email string) ([]*model.Order, error) {
	return r.getOrdersByFilter(ctx, "user_id = $1 AND status = $2 AND phone = $3 AND email = $4", userId, status, phone, email)
}

func (r *repository) UpdateOrderStatus(ctx context.Context, orderId string, status model.OrderStatus) error {
	query := `
		UPDATE orders
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE order_id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, orderId)
	return err
}

func (r *repository) getOrdersByFilter(ctx context.Context, filter string, args ...interface{}) ([]*model.Order, error) {
	query := `
		SELECT o.order_id, o.phone, o.email, o.description, o.status,
			   p.product_id, p.name, p.weight, p.description
		FROM orders o
		JOIN products p ON o.product_id = p.product_id
		WHERE ` + filter

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		var order model.Order
		var product model.Product
		err := rows.Scan(
			&order.OrderId,
			&order.Phone,
			&order.Email,
			&order.Description,
			&order.Status,
			&product.ProductId,
			&product.Name,
			&product.Weigth,
			&product.Description,
		)
		if err != nil {
			return nil, err
		}
		order.Product = product
		orders = append(orders, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}
