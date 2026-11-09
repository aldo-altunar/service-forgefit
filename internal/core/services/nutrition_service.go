package services

import (
	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
	"context"
	"time"

	"github.com/google/uuid"
)

type nutritionService struct{ nutritionRepo ports.NutritionRepository }

func NewNutritionService(repo ports.NutritionRepository) ports.NutritionService {
	return &nutritionService{nutritionRepo: repo}
}

func (s *nutritionService) GetTodayNutrition(ctx context.Context, userID string) (*domain.DailyNutrition, error) {
	today := time.Now().Format("2006-01-02")
	nutrition, err := s.nutritionRepo.GetByDate(ctx, userID, today)
	if err != nil {
		return nil, err
	}
	if nutrition == nil {
		nutrition = &domain.DailyNutrition{ID: uuid.New().String(), UserID: userID, Date: today, WaterGoalMl: 2500, CaloriesGoal: 2200, ProteinGoalG: 180, CarbsGoalG: 250, FatGoalG: 70}
		if err := s.nutritionRepo.UpsertDailyNutrition(ctx, nutrition); err != nil {
			return nil, err
		}
	}
	return nutrition, nil
}

func (s *nutritionService) AddWater(ctx context.Context, userID string, input ports.AddWaterInput) (*domain.DailyNutrition, error) {
	if input.AmountMl != 150 && input.AmountMl != 250 && input.AmountMl != 500 {
		return nil, domain.ErrInvalidNutritionInput
	}
	nutrition, err := s.GetTodayNutrition(ctx, userID)
	if err != nil {
		return nil, err
	}
	today := time.Now().Format("2006-01-02")
	log := &domain.HydrationLog{ID: uuid.New().String(), UserID: userID, Date: today, AmountMl: input.AmountMl, LoggedAt: time.Now()}
	if err := s.nutritionRepo.AddHydrationLog(ctx, log); err != nil {
		return nil, err
	}
	nutrition.WaterIntakeMl += input.AmountMl
	nutrition.UpdatedAt = time.Now()
	if err := s.nutritionRepo.UpsertDailyNutrition(ctx, nutrition); err != nil {
		return nil, err
	}
	return nutrition, nil
}
