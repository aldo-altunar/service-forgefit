package ports

import (
	"context"
	"ForgeFit_Servidor/internal/core/domain"
)

type ActivityRepository interface {
	GetByDate(ctx context.Context, userID string, date string) (*domain.DailyActivity, error)
	Upsert(ctx context.Context, activity *domain.DailyActivity) error
	GetWeeklyActivity(ctx context.Context, userID string) ([]domain.WeeklyActivityDay, error)
}

type SyncStepsInput struct {
	Steps         int     `json:"steps"`
	DistanceKm    float64 `json:"distance_km"`
	Calories      int     `json:"calories"`
	ActiveMinutes int     `json:"active_minutes"`
}

type ActivityService interface {
	GetTodayActivity(ctx context.Context, userID string) (*domain.DailyActivity, error)
	SyncSteps(ctx context.Context, userID string, input SyncStepsInput) (*domain.DailyActivity, error)
	GetWeeklyChart(ctx context.Context, userID string) ([]domain.WeeklyActivityDay, error)
}
