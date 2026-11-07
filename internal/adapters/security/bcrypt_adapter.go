package security

import (
	"golang.org/x/crypto/bcrypt"
	"ForgeFit_Servidor/internal/core/ports"
)

type bcryptAdapter struct {
	cost int
}

func NewBcryptAdapter(cost int) ports.PasswordHasher {
	if cost <= 0 {
		cost = bcrypt.DefaultCost
	}
	return &bcryptAdapter{cost: cost}
}

func (b *bcryptAdapter) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), b.cost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (b *bcryptAdapter) Compare(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
