package orders

import (
	"backend_crm/internal/model"
	"context"
)

type Repository interface {
	Save(ctx context.Context, newOrder *model.NewOrder) error
	GetAll(ctx context.Context) ([]*model.Order, error)
	GetByStatus(ctx context.Context, status model.OrderStatus) ([]*model.Order, error)
	GetByStatusAndPhone(ctx context.Context, status model.OrderStatus, phone string) ([]*model.Order, error)
	GetByStatusAndEmail(ctx context.Context, status model.OrderStatus, email string) ([]*model.Order, error)
	GetByStatusAndPhoneAndEmail(ctx context.Context, status model.OrderStatus, phone string, email string) ([]*model.Order, error)
	GetByUserIdAndStatus(ctx context.Context, userId string, status model.OrderStatus) ([]*model.Order, error)
	GetByUserIdAndStatusAndPhone(ctx context.Context, userId string, status model.OrderStatus, phone string) ([]*model.Order, error)
	GetByUserIdAndStatusAndEmail(ctx context.Context, userId string, status model.OrderStatus, email string) ([]*model.Order, error)
	GetByUserIdAndStatusAndPhoneAndEmail(ctx context.Context, userId string, status model.OrderStatus, phone string, email string) ([]*model.Order, error)
	UpdateOrderStatus(ctx context.Context, orderId string, status model.OrderStatus) error
}
