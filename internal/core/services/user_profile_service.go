package services

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type userProfileService struct {
	profileRepo ports.UserProfileRepository
}

func NewUserProfileService(profileRepo ports.UserProfileRepository) ports.UserProfileService {
	return &userProfileService{profileRepo: profileRepo}
}

func (s *userProfileService) GetProfileByUserID(ctx context.Context, userID string) (*domain.UserProfile, error) {
	profile, err := s.profileRepo.FindByUserID(ctx, userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			newProfile := &domain.UserProfile{
				ID:            uuid.New().String(),
				UserID:        userID,
				WeightKg:      0,
				HeightCm:      0,
				Age:           0,
				TotalWorkouts: 0,
				StreakDays:    0,
				TotalKm:       0,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			if err := s.profileRepo.Create(ctx, newProfile); err != nil {
				return nil, err
			}
			return newProfile, nil
		}
		return nil, err
	}
	return profile, nil
}

func (s *userProfileService) UpdateProfile(ctx context.Context, userID string, input ports.UpdateProfileInput) (*domain.UserProfile, error) {
	profile, err := s.GetProfileByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	if input.WeightKg > 0 {
		profile.WeightKg = input.WeightKg
	}
	if input.HeightCm > 0 {
		profile.HeightCm = input.HeightCm
	}
	if input.Age > 0 {
		profile.Age = input.Age
	}
	profile.UpdatedAt = time.Now()

	if err := s.profileRepo.Update(ctx, profile); err != nil {
		return nil, err
	}

	return profile, nil
}
