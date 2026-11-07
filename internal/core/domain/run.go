package domain

import (
	"errors"
	"time"
)

var (
	ErrRunNotFound     = errors.New("run not found")
	ErrInvalidRunState = errors.New("invalid run state")
	ErrInvalidGPSPoint = errors.New("invalid GPS coordinate")
)

type RunStatus string

const (
	RunStatusInProgress RunStatus = "in_progress"
	RunStatusPaused     RunStatus = "paused"
	RunStatusCompleted  RunStatus = "completed"
	RunStatusCancelled  RunStatus = "cancelled"
)

type Run struct {
	ID              string          `json:"id"`
	UserID          string          `json:"user_id"`
	Status          RunStatus       `json:"status"`
	StartTime       time.Time       `json:"start_time"`
	EndTime         *time.Time      `json:"end_time,omitempty"`
	DurationSeconds int             `json:"duration_seconds"`
	DistanceKm      float64         `json:"distance_km"`
	CurrentSpeedKmh float64         `json:"current_speed_kmh"`
	AvgPaceMinKm    float64         `json:"avg_pace_min_km"`
	CaloriesBurned  int             `json:"calories_burned"`
	AvgHeartRateBpm int             `json:"avg_heart_rate_bpm"`
	ElevationGainM  int             `json:"elevation_gain_m"`
	RoutePoints     []RunRoutePoint `json:"route_points,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type RunRoutePoint struct {
	ID            string    `json:"id"`
	RunID         string    `json:"run_id"`
	Latitude      float64   `json:"latitude"`
	Longitude     float64   `json:"longitude"`
	SpeedKmh      float64   `json:"speed_kmh"`
	AltitudeM     float64   `json:"altitude_m"`
	SequenceOrder int       `json:"sequence_order"`
	Timestamp     time.Time `json:"timestamp"`
}
