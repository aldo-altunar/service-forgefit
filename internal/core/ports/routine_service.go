package ports

import (
	"context"
	"ForgeFit_Servidor/internal/core/domain"
)

type RoutineRepository interface {
	GetTodayRoutine(ctx context.Context, userID string, date string) (*domain.Routine, error)
	GetNextRoutine(ctx context.Context, userID string) (*domain.Routine, error)
	FindByID(ctx context.Context, routineID string) (*domain.Routine, error)
	UpdateExerciseCompletion(ctx context.Context, exerciseID string, isCompleted bool) error
	UpdateRoutineProgress(ctx context.Context, routineID string) error
	CreateRoutine(ctx context.Context, routine *domain.Routine) error
}

type RoutineService interface {
	GetTodayRoutine(ctx context.Context, userID string) (*domain.Routine, error)
	GetNextRoutine(ctx context.Context, userID string) (*domain.Routine, error)
	ToggleExerciseCompletion(ctx context.Context, userID string, exerciseID string, isCompleted bool) (*domain.Routine, error)
}
