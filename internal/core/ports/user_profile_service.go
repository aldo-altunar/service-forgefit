package ports

import (
	"context"
	"ForgeFit_Servidor/internal/core/domain"
)

type UserProfileRepository interface {
	Create(ctx context.Context, profile *domain.UserProfile) error
	FindByUserID(ctx context.Context, userID string) (*domain.UserProfile, error)
	Update(ctx context.Context, profile *domain.UserProfile) error
}

type UpdateProfileInput struct {
	WeightKg float64 `json:"weight_kg"`
	HeightCm float64 `json:"height_cm"`
	Age      int     `json:"age"`
}

type UserProfileService interface {
	GetProfileByUserID(ctx context.Context, userID string) (*domain.UserProfile, error)
	UpdateProfile(ctx context.Context, userID string, input UpdateProfileInput) (*domain.UserProfile, error)
}
