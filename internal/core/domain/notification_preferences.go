package domain

import "time"

type NotificationPreferences struct {
	UserID             string    `json:"user_id"`
	WorkoutReminders   bool      `json:"workout_reminders"`
	HydrationReminders bool      `json:"hydration_reminders"`
	NutritionReminders bool      `json:"nutrition_reminders"`
	WeeklySummary      bool      `json:"weekly_summary"`
	UpdatedAt          time.Time `json:"updated_at"`
}
