package ports

import (
	"context"
	"ForgeFit_Servidor/internal/core/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	DeleteAccount(ctx context.Context, userID string) error
	IsAccountActive(ctx context.Context, userID string) (bool, error)
}
