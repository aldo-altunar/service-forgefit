package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/adapters/handlers/middlewares"
	"ForgeFit_Servidor/internal/core/ports"
)

type ActivityHandler struct {
	activityService ports.ActivityService
}

func NewActivityHandler(activityService ports.ActivityService) *ActivityHandler {
	return &ActivityHandler{activityService: activityService}
}

func (h *ActivityHandler) GetTodayActivity(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	activity, err := h.activityService.GetTodayActivity(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener actividad diaria: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "actividad del día obtenida exitosamente",
		"data":    activity,
	})
}

func (h *ActivityHandler) SyncSteps(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	var input ports.SyncStepsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de pasos inválidos: " + err.Error()})
		return
	}

	activity, err := h.activityService.SyncSteps(c.Request.Context(), userID.(string), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al sincronizar pasos: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "pasos sincronizados exitosamente",
		"data":    activity,
	})
}

func (h *ActivityHandler) GetWeeklyChart(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	chart, err := h.activityService.GetWeeklyChart(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener gráfica semanal: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "gráfica semanal obtenida exitosamente",
		"data":    chart,
	})
}
