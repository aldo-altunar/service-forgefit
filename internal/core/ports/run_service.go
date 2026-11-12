package ports

import (
	"ForgeFit_Servidor/internal/core/domain"
	"context"
)

type RunRepository interface {
	Create(ctx context.Context, run *domain.Run) error
	FindByID(ctx context.Context, runID string) (*domain.Run, error)
	GetActiveRunByUserID(ctx context.Context, userID string) (*domain.Run, error)
	GetLastRunByUserID(ctx context.Context, userID string) (*domain.Run, error)
	ListByUserID(ctx context.Context, userID string, limit int) ([]domain.Run, error)
	DeleteByUserID(ctx context.Context, userID string, runID string) (bool, error)
	Update(ctx context.Context, run *domain.Run) error
	AddRoutePoint(ctx context.Context, point *domain.RunRoutePoint) error
	GetRoutePoints(ctx context.Context, runID string) ([]domain.RunRoutePoint, error)
}

type StartRunInput struct {
	InitialLat float64 `json:"initial_lat"`
	InitialLng float64 `json:"initial_lng"`
}

type AddGPSPointInput struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	SpeedKmh  float64 `json:"speed_kmh"`
	AltitudeM float64 `json:"altitude_m"`
}

type UpdateRunMetricsInput struct {
	DurationSeconds int     `json:"duration_seconds"`
	DistanceKm      float64 `json:"distance_km"`
	CurrentSpeedKmh float64 `json:"current_speed_kmh"`
	AvgPaceMinKm    float64 `json:"avg_pace_min_km"`
	CaloriesBurned  int     `json:"calories_burned"`
	AvgHeartRateBpm int     `json:"avg_heart_rate_bpm"`
	ElevationGainM  int     `json:"elevation_gain_m"`
}

type RunService interface {
	StartRun(ctx context.Context, userID string, input StartRunInput) (*domain.Run, error)
	GetActiveRun(ctx context.Context, userID string) (*domain.Run, error)
	PauseRun(ctx context.Context, userID string, runID string) (*domain.Run, error)
	ResumeRun(ctx context.Context, userID string, runID string) (*domain.Run, error)
	FinishRun(ctx context.Context, userID string, runID string, input UpdateRunMetricsInput) (*domain.Run, error)
	AddGPSPoint(ctx context.Context, userID string, runID string, input AddGPSPointInput) (*domain.RunRoutePoint, error)
	GetLastRun(ctx context.Context, userID string) (*domain.Run, error)
	GetRunHistory(ctx context.Context, userID string) ([]domain.Run, error)
	DeleteRun(ctx context.Context, userID string, runID string) (bool, error)
}
