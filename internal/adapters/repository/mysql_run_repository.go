package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type mysqlRunRepository struct {
	db *sql.DB
}

func NewMySQLRunRepository(db *sql.DB) ports.RunRepository {
	return &mysqlRunRepository{db: db}
}

func (r *mysqlRunRepository) Create(ctx context.Context, run *domain.Run) error {
	query := `
		INSERT INTO runs (id, user_id, status, start_time, duration_seconds, distance_km, current_speed_kmh, avg_pace_min_km, calories_burned, avg_heart_rate_bpm, elevation_gain_m, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, NOW(), NOW())
	`
	_, err := r.db.ExecContext(ctx, query,
		run.ID,
		run.UserID,
		run.Status,
		run.StartTime,
		run.DurationSeconds,
		run.DistanceKm,
		run.CurrentSpeedKmh,
		run.AvgPaceMinKm,
		run.CaloriesBurned,
		run.AvgHeartRateBpm,
		run.ElevationGainM,
	)
	if err != nil {
		return fmt.Errorf("error creating run: %w", err)
	}
	return nil
}

func (r *mysqlRunRepository) FindByID(ctx context.Context, runID string) (*domain.Run, error) {
	query := `
		SELECT id, user_id, status, start_time, end_time, duration_seconds, distance_km, current_speed_kmh, avg_pace_min_km, calories_burned, avg_heart_rate_bpm, elevation_gain_m, created_at, updated_at
		FROM runs WHERE id = ? LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, runID)
	return r.scanRunRow(row)
}

func (r *mysqlRunRepository) GetActiveRunByUserID(ctx context.Context, userID string) (*domain.Run, error) {
	query := `
		SELECT id, user_id, status, start_time, end_time, duration_seconds, distance_km, current_speed_kmh, avg_pace_min_km, calories_burned, avg_heart_rate_bpm, elevation_gain_m, created_at, updated_at
		FROM runs WHERE user_id = ? AND status IN ('in_progress', 'paused') ORDER BY start_time DESC LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, userID)
	return r.scanRunRow(row)
}

func (r *mysqlRunRepository) GetLastRunByUserID(ctx context.Context, userID string) (*domain.Run, error) {
	query := `
		SELECT id, user_id, status, start_time, end_time, duration_seconds, distance_km, current_speed_kmh, avg_pace_min_km, calories_burned, avg_heart_rate_bpm, elevation_gain_m, created_at, updated_at
		FROM runs WHERE user_id = ? AND status = 'completed' ORDER BY start_time DESC LIMIT 1
	`
	row := r.db.QueryRowContext(ctx, query, userID)
	return r.scanRunRow(row)
}

func (r *mysqlRunRepository) ListByUserID(ctx context.Context, userID string, limit int) ([]domain.Run, error) {
	if limit < 1 || limit > 100 {
		limit = 50
	}
	rows, err := r.db.QueryContext(ctx, `SELECT id, user_id, status, start_time, end_time, duration_seconds, distance_km, current_speed_kmh, avg_pace_min_km, calories_burned, avg_heart_rate_bpm, elevation_gain_m, created_at, updated_at FROM runs WHERE user_id = ? ORDER BY start_time DESC LIMIT ?`, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := make([]domain.Run, 0)
	for rows.Next() {
		var item domain.Run
		var endTime sql.NullTime
		if err := rows.Scan(&item.ID, &item.UserID, &item.Status, &item.StartTime, &endTime, &item.DurationSeconds, &item.DistanceKm, &item.CurrentSpeedKmh, &item.AvgPaceMinKm, &item.CaloriesBurned, &item.AvgHeartRateBpm, &item.ElevationGainM, &item.CreatedAt, &item.UpdatedAt); err != nil {
			return nil, err
		}
		if endTime.Valid {
			item.EndTime = &endTime.Time
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *mysqlRunRepository) DeleteByUserID(ctx context.Context, userID string, runID string) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("error starting run deletion: %w", err)
	}
	defer tx.Rollback()
	if _, err = tx.ExecContext(ctx, `DELETE FROM run_route_points WHERE run_id = ? AND EXISTS (SELECT 1 FROM runs WHERE id = ? AND user_id = ?)`, runID, runID, userID); err != nil {
		return false, fmt.Errorf("error deleting run route: %w", err)
	}
	result, err := tx.ExecContext(ctx, `DELETE FROM runs WHERE id = ? AND user_id = ?`, runID, userID)
	if err != nil {
		return false, fmt.Errorf("error deleting run: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("error checking deleted run: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("error committing run deletion: %w", err)
	}
	return count > 0, nil
}

func (r *mysqlRunRepository) Update(ctx context.Context, run *domain.Run) error {
	query := `
		UPDATE runs
		SET status = ?, end_time = ?, duration_seconds = ?, distance_km = ?, current_speed_kmh = ?, avg_pace_min_km = ?, calories_burned = ?, avg_heart_rate_bpm = ?, elevation_gain_m = ?, updated_at = NOW()
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query,
		run.Status,
		run.EndTime,
		run.DurationSeconds,
		run.DistanceKm,
		run.CurrentSpeedKmh,
		run.AvgPaceMinKm,
		run.CaloriesBurned,
		run.AvgHeartRateBpm,
		run.ElevationGainM,
		run.ID,
	)
	if err != nil {
		return fmt.Errorf("error updating run: %w", err)
	}
	return nil
}

func (r *mysqlRunRepository) AddRoutePoint(ctx context.Context, point *domain.RunRoutePoint) error {
	query := `
		INSERT INTO run_route_points (id, run_id, latitude, longitude, speed_kmh, altitude_m, sequence_order, timestamp)
		VALUES (?, ?, ?, ?, ?, ?, ?, NOW())
	`
	_, err := r.db.ExecContext(ctx, query,
		point.ID,
		point.RunID,
		point.Latitude,
		point.Longitude,
		point.SpeedKmh,
		point.AltitudeM,
		point.SequenceOrder,
	)
	if err != nil {
		return fmt.Errorf("error inserting run_route_point: %w", err)
	}
	return nil
}

func (r *mysqlRunRepository) GetRoutePoints(ctx context.Context, runID string) ([]domain.RunRoutePoint, error) {
	query := `
		SELECT id, run_id, latitude, longitude, speed_kmh, altitude_m, sequence_order, timestamp
		FROM run_route_points WHERE run_id = ? ORDER BY sequence_order ASC
	`
	rows, err := r.db.QueryContext(ctx, query, runID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var points []domain.RunRoutePoint
	for rows.Next() {
		var p domain.RunRoutePoint
		if err := rows.Scan(&p.ID, &p.RunID, &p.Latitude, &p.Longitude, &p.SpeedKmh, &p.AltitudeM, &p.SequenceOrder, &p.Timestamp); err == nil {
			points = append(points, p)
		}
	}
	return points, nil
}

func (r *mysqlRunRepository) scanRunRow(row *sql.Row) (*domain.Run, error) {
	var run domain.Run
	var endTime sql.NullTime
	err := row.Scan(
		&run.ID,
		&run.UserID,
		&run.Status,
		&run.StartTime,
		&endTime,
		&run.DurationSeconds,
		&run.DistanceKm,
		&run.CurrentSpeedKmh,
		&run.AvgPaceMinKm,
		&run.CaloriesBurned,
		&run.AvgHeartRateBpm,
		&run.ElevationGainM,
		&run.CreatedAt,
		&run.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	if endTime.Valid {
		run.EndTime = &endTime.Time
	}
	return &run, nil
}
