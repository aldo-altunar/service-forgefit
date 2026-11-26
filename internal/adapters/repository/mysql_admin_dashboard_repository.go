package repository

import (
	"context"
	"database/sql"
	"fmt"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type mysqlAdminDashboardRepository struct {
	db *sql.DB
}

func NewMySQLAdminDashboardRepository(db *sql.DB) ports.AdminDashboardRepository {
	return &mysqlAdminDashboardRepository{db: db}
}

func (r *mysqlAdminDashboardRepository) GetAdminDashboardSummary(ctx context.Context) (*domain.AdminDashboardSummary, error) {
	const summaryQuery = `
		SELECT
			(SELECT COUNT(*) FROM users WHERE user_role = 'user' AND is_deleted = FALSE),
			(SELECT COUNT(*) FROM routines r JOIN users u ON u.id = r.user_id WHERE u.user_role = 'user' AND u.is_deleted = FALSE),
			(SELECT COUNT(*) FROM routine_exercises re JOIN routines r ON r.id = re.routine_id JOIN users u ON u.id = r.user_id WHERE u.user_role = 'user' AND u.is_deleted = FALSE),
			(SELECT COUNT(*) FROM routine_exercises re JOIN routines r ON r.id = re.routine_id JOIN users u ON u.id = r.user_id WHERE u.user_role = 'user' AND u.is_deleted = FALSE AND re.is_completed = TRUE),
			(SELECT COUNT(*) FROM meal_logs ml JOIN users u ON u.id = ml.user_id WHERE u.user_role = 'user' AND u.is_deleted = FALSE),
			(SELECT COALESCE(SUM(dn.water_intake_ml), 0) FROM daily_nutrition dn JOIN users u ON u.id = dn.user_id WHERE u.user_role = 'user' AND u.is_deleted = FALSE AND dn.date = CURRENT_DATE())
	`

	summary := &domain.AdminDashboardSummary{}
	if err := r.db.QueryRowContext(ctx, summaryQuery).Scan(
		&summary.TotalUsers,
		&summary.TotalRoutines,
		&summary.TotalExercises,
		&summary.CompletedExercises,
		&summary.TotalMeals,
		&summary.WaterTodayMl,
	); err != nil {
		return nil, fmt.Errorf("query admin dashboard totals: %w", err)
	}

	const mostActiveQuery = `
		SELECT u.id, u.name, COUNT(*) AS completed_exercises
		FROM routine_exercises re
		JOIN routines r ON r.id = re.routine_id
		JOIN users u ON u.id = r.user_id
		WHERE u.user_role = 'user' AND u.is_deleted = FALSE AND re.is_completed = TRUE
		GROUP BY u.id, u.name
		ORDER BY completed_exercises DESC, u.name ASC
		LIMIT 1
	`

	var topUser domain.AdminDashboardUserStats
	err := r.db.QueryRowContext(ctx, mostActiveQuery).Scan(&topUser.UserID, &topUser.Name, &topUser.CompletedExercises)
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("query most active user: %w", err)
	}
	if err == nil {
		summary.MostActiveUser = &topUser
	}

	return summary, nil
}
