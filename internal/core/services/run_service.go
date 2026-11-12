package services

import (
	"context"
	"math"
	"time"

	"github.com/google/uuid"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type runService struct {
	runRepo ports.RunRepository
}

func NewRunService(runRepo ports.RunRepository) ports.RunService {
	return &runService{runRepo: runRepo}
}

func (s *runService) StartRun(ctx context.Context, userID string, input ports.StartRunInput) (*domain.Run, error) {
	if input.InitialLat < -90 || input.InitialLat > 90 || input.InitialLng < -180 || input.InitialLng > 180 {
		return nil, domain.ErrInvalidGPSPoint
	}
	activeRun, err := s.runRepo.GetActiveRunByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if activeRun != nil {
		return activeRun, nil
	}

	newRun := &domain.Run{
		ID:              uuid.New().String(),
		UserID:          userID,
		Status:          domain.RunStatusInProgress,
		StartTime:       time.Now(),
		DurationSeconds: 0,
		DistanceKm:      0.0,
		CurrentSpeedKmh: 0,
		AvgPaceMinKm:    0,
		CaloriesBurned:  0,
		AvgHeartRateBpm: 0,
		ElevationGainM:  0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := s.runRepo.Create(ctx, newRun); err != nil {
		return nil, err
	}

	if input.InitialLat != 0 || input.InitialLng != 0 {
		initialPoint := &domain.RunRoutePoint{
			ID:            uuid.New().String(),
			RunID:         newRun.ID,
			Latitude:      input.InitialLat,
			Longitude:     input.InitialLng,
			SpeedKmh:      0,
			AltitudeM:     0,
			SequenceOrder: 1,
			Timestamp:     time.Now(),
		}
		if err := s.runRepo.AddRoutePoint(ctx, initialPoint); err != nil {
			return nil, err
		}
	}

	return newRun, nil
}

func (s *runService) GetActiveRun(ctx context.Context, userID string) (*domain.Run, error) {
	run, err := s.runRepo.GetActiveRunByUserID(ctx, userID)
	if err != nil || run == nil {
		return run, err
	}
	if run.Status == domain.RunStatusInProgress {
		run.DurationSeconds = int(time.Since(run.StartTime).Seconds())
	}

	points, _ := s.runRepo.GetRoutePoints(ctx, run.ID)
	run.RoutePoints = points
	return run, nil
}

func (s *runService) PauseRun(ctx context.Context, userID string, runID string) (*domain.Run, error) {
	run, err := s.runRepo.FindByID(ctx, runID)
	if err != nil {
		return nil, err
	}
	if run == nil || run.UserID != userID {
		return nil, domain.ErrRunNotFound
	}
	if run.Status != domain.RunStatusInProgress {
		return nil, domain.ErrInvalidRunState
	}
	run.Status = domain.RunStatusPaused
	run.UpdatedAt = time.Now()
	if err := s.runRepo.Update(ctx, run); err != nil {
		return nil, err
	}
	return run, nil
}

func (s *runService) ResumeRun(ctx context.Context, userID string, runID string) (*domain.Run, error) {
	run, err := s.runRepo.FindByID(ctx, runID)
	if err != nil {
		return nil, err
	}
	if run == nil || run.UserID != userID {
		return nil, domain.ErrRunNotFound
	}
	if run.Status != domain.RunStatusPaused {
		return nil, domain.ErrInvalidRunState
	}
	run.Status = domain.RunStatusInProgress
	run.UpdatedAt = time.Now()
	if err := s.runRepo.Update(ctx, run); err != nil {
		return nil, err
	}
	return run, nil
}

func (s *runService) FinishRun(ctx context.Context, userID string, runID string, input ports.UpdateRunMetricsInput) (*domain.Run, error) {
	run, err := s.runRepo.FindByID(ctx, runID)
	if err != nil {
		return nil, err
	}
	if run == nil || run.UserID != userID {
		return nil, domain.ErrRunNotFound
	}
	if run.Status != domain.RunStatusInProgress && run.Status != domain.RunStatusPaused {
		return nil, domain.ErrInvalidRunState
	}
	now := time.Now()
	run.Status = domain.RunStatusCompleted
	run.EndTime = &now
	if input.DurationSeconds <= 0 {
		run.DurationSeconds = int(now.Sub(run.StartTime).Seconds())
	}
	if input.DurationSeconds > 0 {
		run.DurationSeconds = input.DurationSeconds
	}
	if input.DistanceKm > 0 {
		run.DistanceKm = input.DistanceKm
	}
	if input.CurrentSpeedKmh > 0 {
		run.CurrentSpeedKmh = input.CurrentSpeedKmh
	}
	if input.AvgPaceMinKm > 0 {
		run.AvgPaceMinKm = input.AvgPaceMinKm
	}
	if input.CaloriesBurned > 0 {
		run.CaloriesBurned = input.CaloriesBurned
	}
	if input.AvgHeartRateBpm > 0 {
		run.AvgHeartRateBpm = input.AvgHeartRateBpm
	}
	if input.ElevationGainM > 0 {
		run.ElevationGainM = input.ElevationGainM
	}
	run.UpdatedAt = now

	if err := s.runRepo.Update(ctx, run); err != nil {
		return nil, err
	}
	return run, nil
}

func (s *runService) AddGPSPoint(ctx context.Context, userID string, runID string, input ports.AddGPSPointInput) (*domain.RunRoutePoint, error) {
	run, err := s.runRepo.FindByID(ctx, runID)
	if err != nil || run == nil || run.UserID != userID {
		return nil, domain.ErrRunNotFound
	}
	if run.Status != domain.RunStatusInProgress {
		return nil, domain.ErrInvalidRunState
	}
	if input.Latitude < -90 || input.Latitude > 90 || input.Longitude < -180 || input.Longitude > 180 {
		return nil, domain.ErrInvalidGPSPoint
	}
	points, _ := s.runRepo.GetRoutePoints(ctx, runID)
	seqOrder := len(points) + 1

	point := &domain.RunRoutePoint{
		ID:            uuid.New().String(),
		RunID:         runID,
		Latitude:      input.Latitude,
		Longitude:     input.Longitude,
		SpeedKmh:      input.SpeedKmh,
		AltitudeM:     input.AltitudeM,
		SequenceOrder: seqOrder,
		Timestamp:     time.Now(),
	}

	if err := s.runRepo.AddRoutePoint(ctx, point); err != nil {
		return nil, err
	}
	if len(points) > 0 {
		previous := points[len(points)-1]
		segmentKm := gpsDistanceKm(previous.Latitude, previous.Longitude, point.Latitude, point.Longitude)
		run.DistanceKm += segmentKm
		run.DurationSeconds = int(time.Since(run.StartTime).Seconds())
		if run.DurationSeconds > 0 {
			run.AvgPaceMinKm = float64(run.DurationSeconds) / 60 / run.DistanceKm
		}
		run.CaloriesBurned = int(run.DistanceKm * 60)
		run.ElevationGainM += int(max(0, point.AltitudeM-previous.AltitudeM))
		if input.SpeedKmh > 0 {
			run.CurrentSpeedKmh = input.SpeedKmh
		}
		if err := s.runRepo.Update(ctx, run); err != nil {
			return nil, err
		}
	}
	return point, nil
}

func gpsDistanceKm(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0088
	toRadians := func(degrees float64) float64 { return degrees * 3.141592653589793 / 180 }
	dLat, dLon := toRadians(lat2-lat1), toRadians(lon2-lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(toRadians(lat1))*math.Cos(toRadians(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	return 2 * earthRadiusKm * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
}

func (s *runService) GetLastRun(ctx context.Context, userID string) (*domain.Run, error) {
	run, err := s.runRepo.GetLastRunByUserID(ctx, userID)
	return run, err
}

func (s *runService) GetRunHistory(ctx context.Context, userID string) ([]domain.Run, error) {
	return s.runRepo.ListByUserID(ctx, userID, 50)
}

func (s *runService) DeleteRun(ctx context.Context, userID string, runID string) (bool, error) {
	return s.runRepo.DeleteByUserID(ctx, userID, runID)
}
