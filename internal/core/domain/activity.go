package domain

import "time"

type DailyActivity struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	Date           string    `json:"date"` 
	Steps          int       `json:"steps"`
	StepGoal       int       `json:"step_goal"`
	DistanceKm     float64   `json:"distance_km"`
	CaloriesBurned int       `json:"calories_burned"`
	ActiveMinutes  int       `json:"active_minutes"`
	StreakDays     int       `json:"streak_days"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type WeeklyActivityDay struct {
	DayLabel string  `json:"day_label"`
	Steps    int     `json:"steps"`
	ValuePct float64 `json:"value_pct"`
	IsToday  bool    `json:"is_today"`
}
