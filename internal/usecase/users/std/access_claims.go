package std

import (
	"backend_crm/internal/model"

	"github.com/golang-jwt/jwt/v5"
)

type accessClaims struct {
	UserId   string     `json:"user_id"`
	UserRole model.Role `json:"user_role"`
	jwt.RegisteredClaims
}
