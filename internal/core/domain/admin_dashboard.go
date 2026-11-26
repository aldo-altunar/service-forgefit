package domain

type AdminDashboardSummary struct {
	TotalUsers         int                      `json:"total_users"`
	TotalRoutines      int                      `json:"total_routines"`
	TotalExercises     int                      `json:"total_exercises"`
	CompletedExercises int                      `json:"completed_exercises"`
	TotalMeals         int                      `json:"total_meals"`
	WaterTodayMl       int                      `json:"water_today_ml"`
	MostActiveUser     *AdminDashboardUserStats `json:"most_active_user"`
}

type AdminDashboardUserStats struct {
	UserID             string `json:"user_id"`
	Name               string `json:"name"`
	CompletedExercises int    `json:"completed_exercises"`
}
