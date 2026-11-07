package services

import (
	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
	"context"
	"fmt"
	"strings"
)

type dashboardService struct {
	userRepo     ports.UserRepository
	activitySvc  ports.ActivityService
	runSvc       ports.RunService
	routineSvc   ports.RoutineService
	nutritionSvc ports.NutritionService
}

func NewDashboardService(userRepo ports.UserRepository, activitySvc ports.ActivityService, runSvc ports.RunService, routineSvc ports.RoutineService, nutritionSvc ports.NutritionService) ports.DashboardService {
	return &dashboardService{userRepo: userRepo, activitySvc: activitySvc, runSvc: runSvc, routineSvc: routineSvc, nutritionSvc: nutritionSvc}
}

func (s *dashboardService) GetDashboardSummary(ctx context.Context, userID string) (*domain.DashboardSummary, error) {
	user, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	userName := "Usuario"
	if user != nil && user.Name != "" {
		userName = user.Name
	}
	activity, err := s.activitySvc.GetTodayActivity(ctx, userID)
	if err != nil {
		return nil, err
	}
	if activity == nil {
		activity = &domain.DailyActivity{UserID: userID, StepGoal: 10000}
	}
	nutrition, err := s.nutritionSvc.GetTodayNutrition(ctx, userID)
	if err != nil {
		return nil, err
	}
	if nutrition == nil {
		nutrition = &domain.DailyNutrition{CaloriesGoal: 2200}
	}
	lastRun, err := s.runSvc.GetLastRun(ctx, userID)
	if err != nil {
		return nil, err
	}
	lastRunSummary := "Sin carreras registradas"
	if lastRun != nil {
		lastRunSummary = fmt.Sprintf("Última: %.1f km", lastRun.DistanceKm)
	}
	nextRoutine, err := s.routineSvc.GetNextRoutine(ctx, userID)
	if err != nil {
		return nil, err
	}
	nextRoutineSummary := "Sin rutina programada"
	if nextRoutine != nil {
		nextRoutineSummary = "Próxima: " + nextRoutine.Title
	}
	chart, err := s.activitySvc.GetWeeklyChart(ctx, userID)
	if err != nil {
		return nil, err
	}
	if chart == nil {
		chart = []domain.WeeklyActivityDay{}
	}
	return &domain.DashboardSummary{
		UserGreeting: fmt.Sprintf("Hola,\n%s", userName), UserInitials: getInitials(userName),
		HeaderStats: domain.HeaderStats{TotalCalories: activity.CaloriesBurned, ActiveMinutes: activity.ActiveMinutes, StreakDays: activity.StreakDays},
		TodaySteps:  *activity,
		QuickAccess: domain.QuickAccessSummary{LastRunSummary: lastRunSummary, NextRoutineSummary: nextRoutineSummary,
			NutritionSummary: fmt.Sprintf("%s / %s kcal consumidas", formatInt(nutrition.CaloriesConsumed), formatInt(nutrition.CaloriesGoal)),
			CaloriesConsumed: nutrition.CaloriesConsumed, CaloriesGoal: nutrition.CaloriesGoal},
		WeeklyActivityChart: chart,
	}, nil
}

func getInitials(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "U"
	}
	if len(parts) == 1 {
		return strings.ToUpper(string(parts[0][0]))
	}
	return strings.ToUpper(string(parts[0][0]) + string(parts[1][0]))
}

func formatInt(n int) string { return fmt.Sprintf("%d", n) }
