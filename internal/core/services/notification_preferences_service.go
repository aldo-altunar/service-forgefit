package services

import (
	"context"
	"time"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type notificationPreferencesService struct {
	repository ports.NotificationPreferencesRepository
}

func NewNotificationPreferencesService(repository ports.NotificationPreferencesRepository) ports.NotificationPreferencesService {
	return &notificationPreferencesService{repository: repository}
}

func (s *notificationPreferencesService) Get(ctx context.Context, userID string) (*domain.NotificationPreferences, error) {
	preferences, err := s.repository.FindByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if preferences != nil {
		return preferences, nil
	}
	preferences = &domain.NotificationPreferences{UserID: userID, UpdatedAt: time.Now()}
	if err := s.repository.Upsert(ctx, preferences); err != nil {
		return nil, err
	}
	return preferences, nil
}

func (s *notificationPreferencesService) Update(ctx context.Context, userID string, input ports.UpdateNotificationPreferencesInput) (*domain.NotificationPreferences, error) {
	preferences := &domain.NotificationPreferences{
		UserID:             userID,
		WorkoutReminders:   input.WorkoutReminders,
		HydrationReminders: input.HydrationReminders,
		NutritionReminders: input.NutritionReminders,
		WeeklySummary:      input.WeeklySummary,
		UpdatedAt:          time.Now(),
	}
	if err := s.repository.Upsert(ctx, preferences); err != nil {
		return nil, err
	}
	return preferences, nil
}
