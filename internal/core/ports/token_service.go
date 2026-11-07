package ports

import "ForgeFit_Servidor/internal/core/domain"

type TokenClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

type TokenService interface {
	GenerateToken(user *domain.User) (string, error)
	ValidateToken(tokenString string) (*TokenClaims, error)
}
