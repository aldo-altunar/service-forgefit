package repository

import (
	"context"
	"database/sql"
	"errors"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type mysqlRoutineRepository struct {
	db *sql.DB
}

func NewMySQLRoutineRepository(db *sql.DB) ports.RoutineRepository {
	return &mysqlRoutineRepository{db: db}
}

func (r *mysqlRoutineRepository) GetTodayRoutine(ctx context.Context, userID string, date string) (*domain.Routine, error) {
	query := `
		SELECT id, user_id, title, subtitle, scheduled_date, day_name, status, completed_exercises, total_exercises, created_at, updated_at
		FROM routines WHERE user_id = ? AND scheduled_date = ? LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, userID, date)
	routine, err := r.scanRoutineRow(row)
	if err != nil || routine == nil {
		return nil, err
	}
	exercises, _ := r.getRoutineExercises(ctx, routine.ID)
	routine.Exercises = exercises
	return routine, nil
}

func (r *mysqlRoutineRepository) GetNextRoutine(ctx context.Context, userID string) (*domain.Routine, error) {
	query := `
		SELECT id, user_id, title, subtitle, scheduled_date, day_name, status, completed_exercises, total_exercises, created_at, updated_at
		FROM routines WHERE user_id = ? AND status != 'completed' ORDER BY scheduled_date ASC LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, userID)
	routine, err := r.scanRoutineRow(row)
	if err != nil || routine == nil {
		return nil, err
	}
	exercises, _ := r.getRoutineExercises(ctx, routine.ID)
	routine.Exercises = exercises
	return routine, nil
}

func (r *mysqlRoutineRepository) FindByID(ctx context.Context, routineID string) (*domain.Routine, error) {
	query := `
		SELECT id, user_id, title, subtitle, scheduled_date, day_name, status, completed_exercises, total_exercises, created_at, updated_at
		FROM routines WHERE id = ? LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, routineID)
	routine, err := r.scanRoutineRow(row)
	if err != nil || routine == nil {
		return nil, err
	}
	exercises, _ := r.getRoutineExercises(ctx, routine.ID)
	routine.Exercises = exercises
	return routine, nil
}

func (r *mysqlRoutineRepository) UpdateExerciseCompletion(ctx context.Context, exerciseID string, isCompleted bool) error {
	query := `UPDATE routine_exercises SET is_completed = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, isCompleted, exerciseID)
	return err
}

func (r *mysqlRoutineRepository) UpdateRoutineProgress(ctx context.Context, routineID string) error {
	query := `
		UPDATE routines
		SET completed_exercises = (SELECT COUNT(*) FROM routine_exercises WHERE routine_id = ? AND is_completed = 1),
		    updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, routineID, routineID)
	return err
}

func (r *mysqlRoutineRepository) CreateRoutine(ctx context.Context, routine *domain.Routine) error {
	query := `
		INSERT INTO routines (id, user_id, title, subtitle, scheduled_date, day_name, status, completed_exercises, total_exercises, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, query,
		routine.ID,
		routine.UserID,
		routine.Title,
		routine.Subtitle,
		routine.ScheduledDate,
		routine.DayName,
		routine.Status,
		routine.CompletedExercises,
		routine.TotalExercises,
	)
	if err != nil {
		return err
	}

	for _, ex := range routine.Exercises {
		exQuery := `
			INSERT INTO routine_exercises (id, routine_id, exercise_order, name, sets_summary, weight_kg, is_bodyweight, is_completed, created_at)
			VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())
		`
		if _, err := tx.ExecContext(ctx, exQuery,
			ex.ID,
			routine.ID,
			ex.ExerciseOrder,
			ex.Name,
			ex.SetsSummary,
			ex.WeightKg,
			ex.IsBodyweight,
			ex.IsCompleted,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (r *mysqlRoutineRepository) scanRoutineRow(row *sql.Row) (*domain.Routine, error) {
	var rt domain.Routine
	err := row.Scan(
		&rt.ID,
		&rt.UserID,
		&rt.Title,
		&rt.Subtitle,
		&rt.ScheduledDate,
		&rt.DayName,
		&rt.Status,
		&rt.CompletedExercises,
		&rt.TotalExercises,
		&rt.CreatedAt,
		&rt.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &rt, nil
}

func (r *mysqlRoutineRepository) getRoutineExercises(ctx context.Context, routineID string) ([]domain.RoutineExercise, error) {
	query := `
		SELECT id, routine_id, exercise_order, name, sets_summary, weight_kg, is_bodyweight, is_completed, created_at
		FROM routine_exercises WHERE routine_id = ? ORDER BY exercise_order ASC
	`
	rows, err := r.db.QueryContext(ctx, query, routineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var exercises []domain.RoutineExercise
	for rows.Next() {
		var ex domain.RoutineExercise
		if err := rows.Scan(&ex.ID, &ex.RoutineID, &ex.ExerciseOrder, &ex.Name, &ex.SetsSummary, &ex.WeightKg, &ex.IsBodyweight, &ex.IsCompleted, &ex.CreatedAt); err == nil {
			exercises = append(exercises, ex)
		}
	}
	return exercises, nil
}
