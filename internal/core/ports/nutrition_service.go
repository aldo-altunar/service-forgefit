package ports

import (
	"ForgeFit_Servidor/internal/core/domain"
	"context"
)

type NutritionRepository interface {
	GetByDate(ctx context.Context, userID string, date string) (*domain.DailyNutrition, error)
	UpsertDailyNutrition(ctx context.Context, nutrition *domain.DailyNutrition) error
	AddHydrationLog(ctx context.Context, log *domain.HydrationLog) error
}

type AddWaterInput struct {
	AmountMl int `json:"amount_ml" binding:"required"`
}

type NutritionService interface {
	GetTodayNutrition(ctx context.Context, userID string) (*domain.DailyNutrition, error)
	AddWater(ctx context.Context, userID string, input AddWaterInput) (*domain.DailyNutrition, error)
}
