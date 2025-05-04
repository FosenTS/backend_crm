package postgre

import (
	"backend_crm/internal/model"
	"backend_crm/internal/repository/products"
	"context"
	"database/sql"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) products.Repository {
	return &repository{db: db}
}

func (r *repository) Save(ctx context.Context, product *model.Product) error {
	query := `
		INSERT INTO products (name, weight, description)
		VALUES ($1, $2, $3)
		RETURNING product_id
	`

	err := r.db.QueryRowContext(ctx, query,
		product.Name,
		product.Weigth,
		product.Description,
	).Scan(&product.ProductId)

	return err
}

func (r *repository) GetAll() ([]*model.Product, error) {
	query := `
		SELECT product_id, name, weight, description
		FROM products
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.Product
	for rows.Next() {
		var product model.Product
		err := rows.Scan(
			&product.ProductId,
			&product.Name,
			&product.Weigth,
			&product.Description,
		)
		if err != nil {
			return nil, err
		}
		products = append(products, &product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *repository) GetById(id string) (*model.Product, error) {
	query := `
		SELECT product_id, name, weight, description
		FROM products
		WHERE product_id = $1
	`

	var product model.Product
	err := r.db.QueryRow(query, id).Scan(
		&product.ProductId,
		&product.Name,
		&product.Weigth,
		&product.Description,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &product, nil
}
