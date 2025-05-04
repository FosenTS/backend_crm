package users

import (
	"backend_crm/internal/model"
	"context"
	"errors"
)

var (
	ErrNotFoundUser = errors.New("not found user")
)

type Repository interface {
	GetByUsername(ctx context.Context, username string) (*model.User, error)
	Save(ctx context.Context, register *model.Register) error
}
