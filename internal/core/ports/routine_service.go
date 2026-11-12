package ports

import (
	"ForgeFit_Servidor/internal/core/domain"
	"context"
)

type RoutineRepository interface {
	GetTodayRoutine(ctx context.Context, userID string, date string) (*domain.Routine, error)
	GetNextRoutine(ctx context.Context, userID string) (*domain.Routine, error)
	FindByID(ctx context.Context, routineID string) (*domain.Routine, error)
	UpdateExerciseCompletion(ctx context.Context, exerciseID string, isCompleted bool) error
	UpdateRoutineProgress(ctx context.Context, routineID string) error
	CreateRoutine(ctx context.Context, routine *domain.Routine) error
}

type CreateRoutineInput struct {
	Title         string                       `json:"title"`
	Subtitle      string                       `json:"subtitle"`
	ScheduledDate string                       `json:"scheduled_date"`
	DayName       string                       `json:"day_name"`
	Exercises     []CreateRoutineExerciseInput `json:"exercises"`
}

type CreateRoutineExerciseInput struct {
	Name         string  `json:"name"`
	SetsSummary  string  `json:"sets_summary"`
	WeightKg     float64 `json:"weight_kg"`
	IsBodyweight bool    `json:"is_bodyweight"`
}

type RoutineService interface {
	CreateRoutine(ctx context.Context, userID string, input CreateRoutineInput) (*domain.Routine, error)
	GetTodayRoutine(ctx context.Context, userID string) (*domain.Routine, error)
	GetNextRoutine(ctx context.Context, userID string) (*domain.Routine, error)
	ToggleExerciseCompletion(ctx context.Context, userID string, exerciseID string, isCompleted bool) (*domain.Routine, error)
}
