package ports

import (
	"context"

	"ForgeFit_Servidor/internal/core/domain"
)

type UpdateNotificationPreferencesInput struct {
	WorkoutReminders   bool `json:"workout_reminders"`
	HydrationReminders bool `json:"hydration_reminders"`
	NutritionReminders bool `json:"nutrition_reminders"`
	WeeklySummary      bool `json:"weekly_summary"`
}

type NotificationPreferencesRepository interface {
	FindByUserID(ctx context.Context, userID string) (*domain.NotificationPreferences, error)
	Upsert(ctx context.Context, preferences *domain.NotificationPreferences) error
}

type NotificationPreferencesService interface {
	Get(ctx context.Context, userID string) (*domain.NotificationPreferences, error)
	Update(ctx context.Context, userID string, input UpdateNotificationPreferencesInput) (*domain.NotificationPreferences, error)
}
