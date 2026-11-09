package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type mysqlUserProfileRepository struct {
	db *sql.DB
}

func NewMySQLUserProfileRepository(db *sql.DB) ports.UserProfileRepository {
	return &mysqlUserProfileRepository{db: db}
}

func (r *mysqlUserProfileRepository) Create(ctx context.Context, profile *domain.UserProfile) error {
	query := `
		INSERT INTO user_profiles (id, user_id, weight_kg, height_cm, age, gender, activity_level, total_workouts, streak_days, total_km, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	_, err := r.db.ExecContext(ctx, query,
		profile.ID,
		profile.UserID,
		profile.WeightKg,
		profile.HeightCm,
		profile.Age,
		profile.Gender,
		profile.ActivityLevel,
		profile.TotalWorkouts,
		profile.StreakDays,
		profile.TotalKm,
	)
	if err != nil {
		return fmt.Errorf("error inserting user_profile: %w", err)
	}
	return nil
}

func (r *mysqlUserProfileRepository) FindByUserID(ctx context.Context, userID string) (*domain.UserProfile, error) {
	query := `
		SELECT id, user_id, weight_kg, height_cm, age, gender, activity_level, total_workouts, streak_days, total_km, created_at, updated_at
		FROM user_profiles WHERE user_id = ? LIMIT 1
	`
	var p domain.UserProfile
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&p.ID,
		&p.UserID,
		&p.WeightKg,
		&p.HeightCm,
		&p.Age,
		&p.Gender,
		&p.ActivityLevel,
		&p.TotalWorkouts,
		&p.StreakDays,
		&p.TotalKm,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return &p, nil
}

func (r *mysqlUserProfileRepository) Update(ctx context.Context, profile *domain.UserProfile) error {
	query := `
		UPDATE user_profiles
		SET weight_kg = ?, height_cm = ?, age = ?, gender = ?, activity_level = ?, total_workouts = ?, streak_days = ?, total_km = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		profile.WeightKg,
		profile.HeightCm,
		profile.Age,
		profile.Gender,
		profile.ActivityLevel,
		profile.TotalWorkouts,
		profile.StreakDays,
		profile.TotalKm,
		profile.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating user_profile: %w", err)
	}
	return nil
}
