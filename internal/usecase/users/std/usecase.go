package std

import (
	"backend_crm/internal/model"
	usersRepo "backend_crm/internal/repository/users"
	"backend_crm/internal/usecase/users"
	"context"
	"errors"
	"fmt"
	"hash"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var _ users.Usecase = &usecase{}

type usecase struct {
	users usersRepo.Repository

	passHasher hash.Hash

	accessSecret  []byte
	refreshSecret []byte

	accessExpired  time.Duration
	refreshExpired time.Duration
}

func NewUsecase(
	users usersRepo.Repository,
	accessSecret []byte,
	refreshSecret []byte,
	accessTTL time.Duration,
	refreshTTL time.Duration,
) users.Usecase {
	return &usecase{
		users:          users,
		accessSecret:   accessSecret,
		refreshSecret:  refreshSecret,
		accessExpired:  accessTTL,
		refreshExpired: refreshTTL,
	}
}

func (u *usecase) Register(ctx context.Context, register *model.Register) error {
	b, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("generate from password: %w", err)
	}

	register.Password = string(b)

	if err = u.users.Save(ctx, register); err != nil {
		return fmt.Errorf("save user: %w", err)
	}

	return nil
}

func (u *usecase) Login(ctx context.Context, login *model.Login) (*model.Token, error) {
	user, err := u.users.GetByUsername(ctx, login.Username)
	if err != nil {
		if errors.Is(err, usersRepo.ErrNotFoundUser) {
			return nil, users.ErrNotFoundUser
		}
		return nil, fmt.Errorf("get user: %w", err)
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(login.Password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, users.ErrIncorrectPassword
		}
		return nil, fmt.Errorf("compare hash and password: %w", err)
	}

	tokens, err := u.generateTokens(user.UserId, user.Role)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return tokens, nil
}

func (u *usecase) generateTokens(userId string, userRole model.Role) (*model.Token, error) {
	access, err := u.generateAccess(userId, userRole)
	if err != nil {
		return nil, fmt.Errorf("generate access token: %w", err)
	}

	refresh, err := u.generateRefresh(userId, userRole)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	return &model.Token{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func (u *usecase) generateAccess(userId string, userRole model.Role) (string, error) {
	now := time.Now()
	accessExp := now.Add(u.accessExpired)
	accessClaim := &accessClaims{
		UserId:   userId,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExp),
		},
	}
	access := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaim)

	accessToken, err := access.SignedString(u.accessSecret)
	if err != nil {
		return "", fmt.Errorf("signed string: %w", err)
	}

	return accessToken, nil
}

func (u *usecase) generateRefresh(userId string, userRole model.Role) (string, error) {
	now := time.Now()
	refreshExp := now.Add(u.refreshExpired)
	refreshClaim := &refreshClaims{
		UserId:   userId,
		UserRole: userRole,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(refreshExp),
		},
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaim)

	refreshToken, err := refresh.SignedString(u.refreshSecret)
	if err != nil {
		return "", fmt.Errorf("signed string: %w", err)
	}

	return refreshToken, nil
}

func (u *usecase) CheckAccess(ctx context.Context, accessToken string) (string, model.Role, error) {
	claims := &accessClaims{}
	token, err := jwt.ParseWithClaims(
		accessToken,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return u.accessSecret, nil
		})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", 0, users.ErrExpiredAccessToken
		}
		return "", 0, fmt.Errorf("parse with claims: %w", err)
	}

	if !token.Valid {
		return "", 0, errors.New("invalid token")
	}

	return claims.UserId, claims.UserRole, nil
}

func (u *usecase) RefreshTokens(ctx context.Context, refreshToken string) (*model.Token, error) {
	claims := &refreshClaims{}
	token, err := jwt.ParseWithClaims(
		refreshToken,
		claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return u.refreshSecret, nil
		})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, users.ErrExpiredRefreshToken
		}
		return nil, fmt.Errorf("parse with claims: %w", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	newTokens, err := u.generateTokens(claims.UserId, claims.UserRole)
	if err != nil {
		return nil, fmt.Errorf("generate tokens: %w", err)
	}

	return newTokens, nil
}
