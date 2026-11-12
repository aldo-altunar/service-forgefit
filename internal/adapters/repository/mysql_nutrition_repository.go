package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type mysqlNutritionRepository struct {
	db *sql.DB
}

func NewMySQLNutritionRepository(db *sql.DB) ports.NutritionRepository {
	return &mysqlNutritionRepository{db: db}
}

func (r *mysqlNutritionRepository) GetByDate(ctx context.Context, userID string, date string) (*domain.DailyNutrition, error) {
	query := `
		SELECT id, user_id, date, water_intake_ml, water_goal_ml, calories_consumed, calories_goal, protein_g, protein_goal_g, carbs_g, carbs_goal_g, fat_g, fat_goal_g, created_at, updated_at
		FROM daily_nutrition WHERE user_id = ? AND date = ? LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, userID, date)
	var n domain.DailyNutrition
	err := row.Scan(
		&n.ID,
		&n.UserID,
		&n.Date,
		&n.WaterIntakeMl,
		&n.WaterGoalMl,
		&n.CaloriesConsumed,
		&n.CaloriesGoal,
		&n.ProteinG,
		&n.ProteinGoalG,
		&n.CarbsG,
		&n.CarbsGoalG,
		&n.FatG,
		&n.FatGoalG,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error querying daily_nutrition: %w", err)
	}

	meals, err := r.getMealLogs(ctx, userID, date)
	if err != nil {
		return nil, fmt.Errorf("error querying meal logs: %w", err)
	}
	hydrations, _ := r.getHydrationLogs(ctx, userID, date)
	n.MealLogs = meals
	n.HydrationLogs = hydrations

	return &n, nil
}

func (r *mysqlNutritionRepository) UpsertDailyNutrition(ctx context.Context, n *domain.DailyNutrition) error {
	query := `
		INSERT INTO daily_nutrition (id, user_id, date, water_intake_ml, water_goal_ml, calories_consumed, calories_goal, protein_g, protein_goal_g, carbs_g, carbs_goal_g, fat_g, fat_goal_g, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE
			water_intake_ml = VALUES(water_intake_ml),
			water_goal_ml = VALUES(water_goal_ml),
			calories_consumed = VALUES(calories_consumed),
			calories_goal = VALUES(calories_goal),
			protein_g = VALUES(protein_g),
			protein_goal_g = VALUES(protein_goal_g),
			carbs_g = VALUES(carbs_g),
			carbs_goal_g = VALUES(carbs_goal_g),
			fat_g = VALUES(fat_g),
			fat_goal_g = VALUES(fat_goal_g),
			updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query,
		n.ID,
		n.UserID,
		n.Date,
		n.WaterIntakeMl,
		n.WaterGoalMl,
		n.CaloriesConsumed,
		n.CaloriesGoal,
		n.ProteinG,
		n.ProteinGoalG,
		n.CarbsG,
		n.CarbsGoalG,
		n.FatG,
		n.FatGoalG,
	)
	if err != nil {
		return fmt.Errorf("error upserting daily_nutrition: %w", err)
	}
	return nil
}

func (r *mysqlNutritionRepository) AddHydrationLog(ctx context.Context, log *domain.HydrationLog) error {
	query := `
		INSERT INTO hydration_logs (id, user_id, date, amount_ml, logged_at)
		VALUES (?, ?, ?, ?, NOW())
	`
	_, err := r.db.ExecContext(ctx, query, log.ID, log.UserID, log.Date, log.AmountMl)
	return err
}

func (r *mysqlNutritionRepository) AddMealLog(ctx context.Context, log *domain.MealLog) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error starting meal transaction: %w", err)
	}
	defer tx.Rollback()
	_, err = tx.ExecContext(ctx, `INSERT INTO meal_logs (id, user_id, date, meal_type, time_str, name, description, calories, created_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, NOW())`,
		log.ID, log.UserID, log.Date, log.MealType, log.TimeStr, log.Name, log.Description, log.Calories)
	if err != nil {
		return fmt.Errorf("error inserting meal log: %w", err)
	}
	_, err = tx.ExecContext(ctx, `UPDATE daily_nutrition SET calories_consumed = calories_consumed + ?, updated_at = NOW() WHERE user_id = ? AND date = ?`, log.Calories, log.UserID, log.Date)
	if err != nil {
		return fmt.Errorf("error updating daily calories: %w", err)
	}
	return tx.Commit()
}

func (r *mysqlNutritionRepository) getMealLogs(ctx context.Context, userID string, date string) ([]domain.MealLog, error) {
	query := `
		SELECT id, user_id, date, meal_type, time_str, name, description, calories, created_at
		FROM meal_logs WHERE user_id = ? AND date = ? ORDER BY time_str ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var meals []domain.MealLog
	for rows.Next() {
		var m domain.MealLog
		if err := rows.Scan(&m.ID, &m.UserID, &m.Date, &m.MealType, &m.TimeStr, &m.Name, &m.Description, &m.Calories, &m.CreatedAt); err != nil {
			return nil, err
		}
		meals = append(meals, m)
	}
	return meals, rows.Err()
}

func (r *mysqlNutritionRepository) getHydrationLogs(ctx context.Context, userID string, date string) ([]domain.HydrationLog, error) {
	query := `
		SELECT id, user_id, date, amount_ml, logged_at
		FROM hydration_logs WHERE user_id = ? AND date = ? ORDER BY logged_at ASC
	`
	rows, err := r.db.QueryContext(ctx, query, userID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []domain.HydrationLog
	for rows.Next() {
		var h domain.HydrationLog
		if err := rows.Scan(&h.ID, &h.UserID, &h.Date, &h.AmountMl, &h.LoggedAt); err == nil {
			logs = append(logs, h)
		}
	}
	return logs, nil
}
