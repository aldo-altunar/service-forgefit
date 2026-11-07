package domain

type DashboardSummary struct {
	UserGreeting        string              `json:"user_greeting"`         
	UserInitials        string              `json:"user_initials"`         
	HeaderStats         HeaderStats         `json:"header_stats"`          
	TodaySteps          DailyActivity       `json:"today_steps"`           
	QuickAccess         QuickAccessSummary  `json:"quick_access"`          
	WeeklyActivityChart []WeeklyActivityDay `json:"weekly_activity_chart"` 
}

type HeaderStats struct {
	TotalCalories int `json:"total_calories"` 
	ActiveMinutes int `json:"active_minutes"` 
	StreakDays    int `json:"streak_days"`    
}

type QuickAccessSummary struct {
	LastRunSummary      string `json:"last_run_summary"`      
	NextRoutineSummary  string `json:"next_routine_summary"`  
	NutritionSummary    string `json:"nutrition_summary"`    
	CaloriesConsumed    int    `json:"calories_consumed"`
	CaloriesGoal        int    `json:"calories_goal"`
}
