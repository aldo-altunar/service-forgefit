package services

import (
	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

type routineService struct{ routineRepo ports.RoutineRepository }

func NewRoutineService(repo ports.RoutineRepository) ports.RoutineService {
	return &routineService{routineRepo: repo}
}

func (s *routineService) GetTodayRoutine(ctx context.Context, userID string) (*domain.Routine, error) {
	return s.routineRepo.GetTodayRoutine(ctx, userID, time.Now().Format("2006-01-02"))
}

func (s *routineService) GetNextRoutine(ctx context.Context, userID string) (*domain.Routine, error) {
	return s.routineRepo.GetNextRoutine(ctx, userID)
}

func (s *routineService) CreateRoutine(ctx context.Context, userID string, input ports.CreateRoutineInput) (*domain.Routine, error) {
	date, err := time.Parse("2006-01-02", input.ScheduledDate)
	if strings.TrimSpace(input.Title) == "" || err != nil || len(input.Exercises) == 0 {
		return nil, domain.ErrInvalidRoutineInput
	}
	exercises := make([]domain.RoutineExercise, 0, len(input.Exercises))
	for index, exercise := range input.Exercises {
		if strings.TrimSpace(exercise.Name) == "" || strings.TrimSpace(exercise.SetsSummary) == "" || exercise.WeightKg < 0 {
			return nil, domain.ErrInvalidRoutineInput
		}
		exercises = append(exercises, domain.RoutineExercise{
			ID: uuid.New().String(), ExerciseOrder: index + 1,
			Name: strings.TrimSpace(exercise.Name), SetsSummary: strings.TrimSpace(exercise.SetsSummary),
			WeightKg: exercise.WeightKg, IsBodyweight: exercise.IsBodyweight,
		})
	}
	now := time.Now()
	routine := &domain.Routine{
		ID: uuid.New().String(), UserID: userID, Title: strings.TrimSpace(input.Title),
		Subtitle: strings.TrimSpace(input.Subtitle), ScheduledDate: date.Format("2006-01-02"),
		DayName: strings.TrimSpace(input.DayName), Status: domain.RoutineStatusPending,
		TotalExercises: len(exercises), Exercises: exercises, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.routineRepo.CreateRoutine(ctx, routine); err != nil {
		return nil, err
	}
	for i := range routine.Exercises {
		routine.Exercises[i].RoutineID = routine.ID
		routine.Exercises[i].CreatedAt = now
	}
	return routine, nil
}

func (s *routineService) ToggleExerciseCompletion(ctx context.Context, userID, exerciseID string, isCompleted bool) (*domain.Routine, error) {
	routine, err := s.GetTodayRoutine(ctx, userID)
	if err != nil {
		return nil, err
	}
	if routine == nil {
		return nil, errors.New("no routine scheduled for today")
	}
	belongsToRoutine := false
	for _, exercise := range routine.Exercises {
		if exercise.ID == exerciseID {
			belongsToRoutine = true
			break
		}
	}
	if !belongsToRoutine {
		return nil, errors.New("exercise not found in user's routine")
	}
	if err := s.routineRepo.UpdateExerciseCompletion(ctx, exerciseID, isCompleted); err != nil {
		return nil, err
	}
	routine, err = s.GetTodayRoutine(ctx, userID)
	if err != nil {
		return nil, err
	}
	if routine == nil {
		return nil, errors.New("routine disappeared after exercise update")
	}
	count := 0
	for _, exercise := range routine.Exercises {
		if exercise.IsCompleted {
			count++
		}
	}
	routine.CompletedExercises = count
	if count == routine.TotalExercises && count > 0 {
		routine.Status = domain.RoutineStatusCompleted
	} else if count > 0 {
		routine.Status = domain.RoutineStatusInProgress
	} else {
		routine.Status = domain.RoutineStatusPending
	}
	if err := s.routineRepo.UpdateRoutineProgress(ctx, routine.ID); err != nil {
		return nil, err
	}
	return routine, nil
}
