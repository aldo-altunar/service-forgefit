package domain

import (
	"errors"
	"time"
)

type AuthProvider string

const (
	ProviderLocal AuthProvider = "local"
)

type User struct {
	ID           string       `json:"id"`
	Name         string       `json:"name"`
	Email        string       `json:"email"`
	PasswordHash *string      `json:"-"`
	AuthProvider AuthProvider `json:"auth_provider"`
	AvatarURL    *string      `json:"avatar_url,omitempty"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

type UserProfile struct {
	ID            string    `json:"id"`
	UserID        string    `json:"user_id"`
	WeightKg      float64   `json:"weight_kg"`
	HeightCm      float64   `json:"height_cm"`
	Age           int       `json:"age"`
	Gender        string    `json:"gender"`
	ActivityLevel string    `json:"activity_level"`
	TotalWorkouts int       `json:"total_workouts"`
	StreakDays    int       `json:"streak_days"`
	TotalKm       float64   `json:"total_km"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

var (
	ErrUserAlreadyExists  = errors.New("a user with this email already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidToken       = errors.New("invalid or expired token")
	ErrUnauthorized       = errors.New("unauthorized access")
)
