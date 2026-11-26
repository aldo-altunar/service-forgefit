package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/core/ports"
)

type AdminDashboardHandler struct {
	service ports.AdminDashboardService
}

func NewAdminDashboardHandler(service ports.AdminDashboardService) *AdminDashboardHandler {
	return &AdminDashboardHandler{service: service}
}

func (h *AdminDashboardHandler) GetSummary(c *gin.Context) {
	summary, err := h.service.GetAdminDashboardSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener el resumen administrativo"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "resumen administrativo obtenido exitosamente",
		"data":    summary,
	})
}
