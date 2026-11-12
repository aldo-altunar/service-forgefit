package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/adapters/handlers"
	"ForgeFit_Servidor/internal/adapters/handlers/middlewares"
	"ForgeFit_Servidor/internal/adapters/repository"
	"ForgeFit_Servidor/internal/adapters/security"
	"ForgeFit_Servidor/internal/core/services"
	"ForgeFit_Servidor/pkg/config"
	"ForgeFit_Servidor/pkg/database"
)

func main() {
	cfg := config.LoadConfig()
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatalf("Database is unavailable: %v", err)
	}
	defer db.Close()

	userRepo := repository.NewMySQLUserRepository(db)
	userProfileRepo := repository.NewMySQLUserProfileRepository(db)
	activityRepo := repository.NewMySQLActivityRepository(db)
	runRepo := repository.NewMySQLRunRepository(db)
	routineRepo := repository.NewMySQLRoutineRepository(db)
	nutritionRepo := repository.NewMySQLNutritionRepository(db)

	bcryptHasher := security.NewBcryptAdapter(10)
	jwtService := security.NewJWTAdapter(cfg.JWTSecret, cfg.JWTExpirationHours)

	userService := services.NewUserService(userRepo, bcryptHasher, jwtService)
	userProfileService := services.NewUserProfileService(userProfileRepo)
	activityService := services.NewActivityService(activityRepo)
	runService := services.NewRunService(runRepo)
	routineService := services.NewRoutineService(routineRepo)
	nutritionService := services.NewNutritionService(nutritionRepo)
	dashboardService := services.NewDashboardService(userRepo, activityService, runService, routineService, nutritionService)

	userHandler := handlers.NewUserHandler(userService)
	profileHandler := handlers.NewProfileHandler(userProfileService)
	activityHandler := handlers.NewActivityHandler(activityService)
	runHandler := handlers.NewRunHandler(runService)
	routineHandler := handlers.NewRoutineHandler(routineService)
	nutritionHandler := handlers.NewNutritionHandler(nutritionService)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)
	notificationRepo := repository.NewMySQLNotificationPreferencesRepository(db)
	notificationService := services.NewNotificationPreferencesService(notificationRepo)
	notificationHandler := handlers.NewNotificationPreferencesHandler(notificationService)

	router := gin.Default()
	router.Use(middlewares.CORSMiddleware())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "UP", "service": "ForgeFit API"})
	})
	router.GET("/health/ready", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 3*time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{"status": "DOWN"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "READY"})
	})

	api := router.Group("/api/v1")
	{
		auth := api.Group("/auth")
		auth.POST("/register", userHandler.Register)
		auth.POST("/login", userHandler.Login)

		protected := api.Group("")
		protected.Use(middlewares.AuthMiddleware(jwtService, userRepo))
		protected.GET("/dashboard", dashboardHandler.GetDashboardSummary)
		protected.GET("/users/profile", userHandler.GetProfile)
		protected.GET("/users/settings", userHandler.GetAccountSettings)
		protected.PUT("/users/settings", userHandler.UpdateAccountSettings)
		protected.PUT("/users/password", userHandler.ChangePasswordDisabled)
		protected.DELETE("/users/account", userHandler.DeleteAccount)
		protected.GET("/users/notification-preferences", notificationHandler.GetPreferences)
		protected.PUT("/users/notification-preferences", notificationHandler.UpdatePreferences)
		protected.GET("/users/profile/details", profileHandler.GetProfileDetails)
		protected.PUT("/users/profile", profileHandler.UpdateProfile)
		protected.GET("/activity/today", activityHandler.GetTodayActivity)
		protected.POST("/activity/steps", activityHandler.SyncSteps)
		protected.GET("/activity/weekly", activityHandler.GetWeeklyChart)
		protected.POST("/runs/start", runHandler.StartRun)
		protected.GET("/runs/active", runHandler.GetActiveRun)
		protected.GET("/runs/history", runHandler.GetRunHistory)
		protected.DELETE("/runs/:id", runHandler.DeleteRun)
		protected.POST("/runs/:id/pause", runHandler.PauseRun)
		protected.POST("/runs/:id/resume", runHandler.ResumeRun)
		protected.POST("/runs/:id/finish", runHandler.FinishRun)
		protected.POST("/runs/:id/gps", runHandler.AddGPSPoint)
		protected.GET("/routines/today", routineHandler.GetTodayRoutine)
		protected.POST("/routines/exercises/:id/complete", routineHandler.ToggleExercise)
		protected.GET("/nutrition/today", nutritionHandler.GetTodayNutrition)
		protected.POST("/nutrition/water", nutritionHandler.AddWater)
	}

	log.Printf("API listening on port %s", cfg.Port)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
