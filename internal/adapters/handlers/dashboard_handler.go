package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/adapters/handlers/middlewares"
	"ForgeFit_Servidor/internal/core/ports"
)

type DashboardHandler struct {
	dashboardService ports.DashboardService
}

func NewDashboardHandler(dashboardService ports.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) GetDashboardSummary(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	summary, err := h.dashboardService.GetDashboardSummary(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener resumen del dashboard: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "resumen del dashboard obtenido exitosamente",
		"data":    summary,
	})
}
