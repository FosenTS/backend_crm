package users

import (
	"backend_crm/internal/model"
	"context"
	"errors"
)

var (
	ErrNotFoundUser        = errors.New("not found user")
	ErrIncorrectPassword   = errors.New("incorrect password")
	ErrExpiredAccessToken  = errors.New("expired access token")
	ErrExpiredRefreshToken = errors.New("expired refresh token")
)

type Usecase interface {
	CheckAccess(ctx context.Context, accessToken string) (string, model.Role, error)
	RefreshTokens(ctx context.Context, refreshToken string) (*model.Token, error)
	Login(ctx context.Context, login *model.Login) (*model.Token, error)
	Register(ctx context.Context, register *model.Register) error
}
