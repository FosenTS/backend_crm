package products

import (
	"backend_crm/internal/model"
	"context"
)

type Repository interface {
	Save(ctx context.Context, product *model.Product) error
	GetAll() ([]*model.Product, error)
	GetById(id string) (*model.Product, error)
}
