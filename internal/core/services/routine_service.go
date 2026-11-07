package services

import (
	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
	"context"
	"errors"
	"time"
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
