package services

import (
	"context"
	"time"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type activityService struct {
	activityRepo ports.ActivityRepository
}

func NewActivityService(activityRepo ports.ActivityRepository) ports.ActivityService {
	return &activityService{activityRepo: activityRepo}
}

func (s *activityService) GetTodayActivity(ctx context.Context, userID string) (*domain.DailyActivity, error) {
	today := time.Now().Format("2006-01-02")
	activity, err := s.activityRepo.GetByDate(ctx, userID, today)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		activity = &domain.DailyActivity{
			UserID:   userID,
			Date:     today,
			StepGoal: 10000,
		}
		if err := s.activityRepo.Upsert(ctx, activity); err != nil {
			return nil, err
		}
	}
	return activity, nil
}

func (s *activityService) SyncSteps(ctx context.Context, userID string, input ports.SyncStepsInput) (*domain.DailyActivity, error) {
	activity, err := s.GetTodayActivity(ctx, userID)
	if err != nil {
		return nil, err
	}
	activity.Steps = input.Steps
	if input.DistanceKm > 0 {
		activity.DistanceKm = input.DistanceKm
	}
	if input.Calories > 0 {
		activity.CaloriesBurned = input.Calories
	}
	if input.ActiveMinutes > 0 {
		activity.ActiveMinutes = input.ActiveMinutes
	}
	activity.UpdatedAt = time.Now()

	if err := s.activityRepo.Upsert(ctx, activity); err != nil {
		return nil, err
	}
	return activity, nil
}

func (s *activityService) GetWeeklyChart(ctx context.Context, userID string) ([]domain.WeeklyActivityDay, error) {
	chart, err := s.activityRepo.GetWeeklyActivity(ctx, userID)
	if err != nil {
		return nil, err
	}
	if chart == nil {
		chart = []domain.WeeklyActivityDay{}
	}
	return chart, nil
}
