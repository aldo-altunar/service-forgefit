package domain

import "time"

type RoutineStatus string

const (
	RoutineStatusPending    RoutineStatus = "pending"
	RoutineStatusInProgress RoutineStatus = "in_progress"
	RoutineStatusCompleted  RoutineStatus = "completed"
)

type Routine struct {
	ID                 string            `json:"id"`
	UserID             string            `json:"user_id"`
	Title              string            `json:"title"`        
	Subtitle           string            `json:"subtitle"`     
	ScheduledDate      string            `json:"scheduled_date"` 
	DayName            string            `json:"day_name"`
	Status             RoutineStatus     `json:"status"`
	CompletedExercises int               `json:"completed_exercises"`
	TotalExercises     int               `json:"total_exercises"`
	Exercises          []RoutineExercise `json:"exercises,omitempty"`
	CreatedAt          time.Time         `json:"created_at"`
	UpdatedAt          time.Time         `json:"updated_at"`
}

type RoutineExercise struct {
	ID            string  `json:"id"`
	RoutineID     string  `json:"routine_id"`
	ExerciseOrder int     `json:"exercise_order"` 
	Name          string  `json:"name"`           
	SetsSummary   string  `json:"sets_summary"`   
	WeightKg      float64 `json:"weight_kg"`      
	IsBodyweight  bool    `json:"is_bodyweight"`  
	IsCompleted   bool    `json:"is_completed"`
	CreatedAt     time.Time `json:"created_at"`
}
