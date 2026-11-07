package domain

import (
	"errors"
	"time"
)

var ErrInvalidNutritionInput = errors.New("invalid nutrition input")

type MealType string

const (
	MealTypeDesayuno MealType = "desayuno"
	MealTypeAlmuerzo MealType = "almuerzo"
	MealTypeMerienda MealType = "merienda"
	MealTypeCena     MealType = "cena"
)

type DailyNutrition struct {
	ID               string         `json:"id"`
	UserID           string         `json:"user_id"`
	Date             string         `json:"date"`
	WaterIntakeMl    int            `json:"water_intake_ml"`
	WaterGoalMl      int            `json:"water_goal_ml"`
	CaloriesConsumed int            `json:"calories_consumed"`
	CaloriesGoal     int            `json:"calories_goal"`
	ProteinG         int            `json:"protein_g"`
	ProteinGoalG     int            `json:"protein_goal_g"`
	CarbsG           int            `json:"carbs_g"`
	CarbsGoalG       int            `json:"carbs_goal_g"`
	FatG             int            `json:"fat_g"`
	FatGoalG         int            `json:"fat_goal_g"`
	MealLogs         []MealLog      `json:"meal_logs,omitempty"`
	HydrationLogs    []HydrationLog `json:"hydration_logs,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

type HydrationLog struct {
	ID       string    `json:"id"`
	UserID   string    `json:"user_id"`
	Date     string    `json:"date"`
	AmountMl int       `json:"amount_ml"`
	LoggedAt time.Time `json:"logged_at"`
}

type MealLog struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Date        string    `json:"date"`
	MealType    MealType  `json:"meal_type"`
	TimeStr     string    `json:"time_str"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Calories    int       `json:"calories"`
	CreatedAt   time.Time `json:"created_at"`
}
