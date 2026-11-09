package repository

import (
	"context"
	"database/sql"
	"errors"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type mysqlNotificationPreferencesRepository struct {
	db *sql.DB
}

func NewMySQLNotificationPreferencesRepository(db *sql.DB) ports.NotificationPreferencesRepository {
	return &mysqlNotificationPreferencesRepository{db: db}
}

func (r *mysqlNotificationPreferencesRepository) FindByUserID(ctx context.Context, userID string) (*domain.NotificationPreferences, error) {
	preferences := &domain.NotificationPreferences{}
	err := r.db.QueryRowContext(ctx, `
		SELECT user_id, workout_reminders, hydration_reminders, nutrition_reminders, weekly_summary, updated_at
		FROM notification_preferences WHERE user_id = ? LIMIT 1
	`, userID).Scan(
		&preferences.UserID,
		&preferences.WorkoutReminders,
		&preferences.HydrationReminders,
		&preferences.NutritionReminders,
		&preferences.WeeklySummary,
		&preferences.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return preferences, nil
}

func (r *mysqlNotificationPreferencesRepository) Upsert(ctx context.Context, preferences *domain.NotificationPreferences) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO notification_preferences (user_id, workout_reminders, hydration_reminders, nutrition_reminders, weekly_summary, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW())
		ON DUPLICATE KEY UPDATE workout_reminders = VALUES(workout_reminders),
			hydration_reminders = VALUES(hydration_reminders), nutrition_reminders = VALUES(nutrition_reminders),
			weekly_summary = VALUES(weekly_summary), updated_at = NOW()
	`, preferences.UserID, preferences.WorkoutReminders, preferences.HydrationReminders, preferences.NutritionReminders, preferences.WeeklySummary)
	return err
}
