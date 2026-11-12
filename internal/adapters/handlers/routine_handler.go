package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/adapters/handlers/middlewares"
	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type RoutineHandler struct {
	routineService ports.RoutineService
}

func (h *RoutineHandler) CreateRoutine(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	var input ports.CreateRoutineInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de rutina inválidos"})
		return
	}
	routine, err := h.routineService.CreateRoutine(c.Request.Context(), userID.(string), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidRoutineInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al crear rutina: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "rutina creada", "data": routine})
}

func NewRoutineHandler(routineService ports.RoutineService) *RoutineHandler {
	return &RoutineHandler{routineService: routineService}
}

func (h *RoutineHandler) GetTodayRoutine(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	routine, err := h.routineService.GetTodayRoutine(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener rutina de hoy: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "rutina del día obtenida",
		"data":    routine,
	})
}

func (h *RoutineHandler) ToggleExercise(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	exerciseID := c.Param("id")

	type toggleInput struct {
		IsCompleted bool `json:"is_completed"`
	}
	var input toggleInput
	if err := c.ShouldBindJSON(&input); err != nil {
		input.IsCompleted = true
	}

	routine, err := h.routineService.ToggleExerciseCompletion(c.Request.Context(), userID.(string), exerciseID, input.IsCompleted)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar ejercicio: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "estado de ejercicio actualizado",
		"data":    routine,
	})
}
