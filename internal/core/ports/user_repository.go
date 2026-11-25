package ports

import (
	"ForgeFit_Servidor/internal/core/domain"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	DeleteAccount(ctx context.Context, userID string) error
	IsAccountActive(ctx context.Context, userID string) (bool, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)
	CreateFirstAdmin(ctx context.Context, user *domain.User) error
}
