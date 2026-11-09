package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/adapters/handlers/middlewares"
	"ForgeFit_Servidor/internal/core/ports"
)

type NotificationPreferencesHandler struct {
	service ports.NotificationPreferencesService
}

func NewNotificationPreferencesHandler(service ports.NotificationPreferencesService) *NotificationPreferencesHandler {
	return &NotificationPreferencesHandler{service: service}
}

func (h *NotificationPreferencesHandler) GetPreferences(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	preferences, err := h.service.Get(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener preferencias: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "preferencias de notificación obtenidas", "data": preferences})
}

func (h *NotificationPreferencesHandler) UpdatePreferences(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	var input ports.UpdateNotificationPreferencesInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "preferencias inválidas: " + err.Error()})
		return
	}
	preferences, err := h.service.Update(c.Request.Context(), userID.(string), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar preferencias: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "preferencias de notificación actualizadas", "data": preferences})
}
