package security

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type jwtAdapter struct {
	secretKey       []byte
	expirationHours int
}

type jwtCustomClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewJWTAdapter(secretKey string, expirationHours int) ports.TokenService {
	return &jwtAdapter{
		secretKey:       []byte(secretKey),
		expirationHours: expirationHours,
	}
}

func (j *jwtAdapter) GenerateToken(user *domain.User) (string, error) {
	claims := &jwtCustomClaims{
		UserID: user.ID,
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expirationHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

func (j *jwtAdapter) ValidateToken(tokenString string) (*ports.TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
	})

	if err != nil || !token.Valid {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := token.Claims.(*jwtCustomClaims)
	if !ok {
		return nil, domain.ErrInvalidToken
	}

	return &ports.TokenClaims{
		UserID: claims.UserID,
		Email:  claims.Email,
	}, nil
}
