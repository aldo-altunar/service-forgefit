package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/adapters/handlers/middlewares"
	"ForgeFit_Servidor/internal/core/ports"
)

type ProfileHandler struct {
	profileService ports.UserProfileService
}

func NewProfileHandler(profileService ports.UserProfileService) *ProfileHandler {
	return &ProfileHandler{profileService: profileService}
}

func (h *ProfileHandler) GetProfileDetails(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	profile, err := h.profileService.GetProfileByUserID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener detalles del perfil: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "perfil obtenido exitosamente",
		"data":    profile,
	})
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	var input ports.UpdateProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de perfil inválidos: " + err.Error()})
		return
	}

	updated, err := h.profileService.UpdateProfile(c.Request.Context(), userID.(string), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar perfil: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "perfil actualizado exitosamente",
		"data":    updated,
	})
}
