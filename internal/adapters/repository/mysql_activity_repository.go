package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type mysqlActivityRepository struct {
	db *sql.DB
}

func NewMySQLActivityRepository(db *sql.DB) ports.ActivityRepository {
	return &mysqlActivityRepository{db: db}
}

func (r *mysqlActivityRepository) GetByDate(ctx context.Context, userID string, date string) (*domain.DailyActivity, error) {
	query := `
		SELECT id, user_id, date, steps, step_goal, distance_km, calories_burned, active_minutes, streak_days, created_at, updated_at
		FROM daily_activities WHERE user_id = ? AND date = ? LIMIT 1
	`
	var a domain.DailyActivity
	err := r.db.QueryRowContext(ctx, query, userID, date).Scan(
		&a.ID, &a.UserID, &a.Date, &a.Steps, &a.StepGoal, &a.DistanceKm,
		&a.CaloriesBurned, &a.ActiveMinutes, &a.StreakDays, &a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query daily activity: %w", err)
	}
	return &a, nil
}

func (r *mysqlActivityRepository) Upsert(ctx context.Context, activity *domain.DailyActivity) error {
	query := `
		INSERT INTO daily_activities (id, user_id, date, steps, step_goal, distance_km, calories_burned, active_minutes, streak_days, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
		ON DUPLICATE KEY UPDATE steps = VALUES(steps), step_goal = VALUES(step_goal),
			distance_km = VALUES(distance_km), calories_burned = VALUES(calories_burned),
			active_minutes = VALUES(active_minutes), streak_days = VALUES(streak_days), updated_at = NOW()
	`
	_, err := r.db.ExecContext(ctx, query, activity.ID, activity.UserID, activity.Date,
		activity.Steps, activity.StepGoal, activity.DistanceKm, activity.CaloriesBurned,
		activity.ActiveMinutes, activity.StreakDays)
	if err != nil {
		return fmt.Errorf("upsert daily activity: %w", err)
	}
	return nil
}

func (r *mysqlActivityRepository) GetWeeklyActivity(ctx context.Context, userID string) ([]domain.WeeklyActivityDay, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT date, steps FROM daily_activities
		WHERE user_id = ? AND date >= DATE_SUB(CURDATE(), INTERVAL 6 DAY) AND date <= CURDATE()
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stepsByDate := make(map[string]int)
	for rows.Next() {
		var date time.Time
		var steps int
		if err := rows.Scan(&date, &steps); err != nil {
			return nil, err
		}
		stepsByDate[date.Format("2006-01-02")] = steps
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	labels := [...]string{"D", "L", "M", "X", "J", "V", "S"}
	today := time.Now()
	result := make([]domain.WeeklyActivityDay, 0, 7)
	for offset := 6; offset >= 0; offset-- {
		day := today.AddDate(0, 0, -offset)
		steps := stepsByDate[day.Format("2006-01-02")]
		pct := float64(steps) / 10000
		if pct > 1 {
			pct = 1
		}
		result = append(result, domain.WeeklyActivityDay{
			DayLabel: labels[day.Weekday()], Steps: steps, ValuePct: pct,
			IsToday: offset == 0,
		})
	}
	return result, nil
}
