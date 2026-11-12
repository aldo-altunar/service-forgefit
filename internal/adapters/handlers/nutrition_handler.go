package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/adapters/handlers/middlewares"
	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type NutritionHandler struct {
	nutritionService ports.NutritionService
}

func NewNutritionHandler(nutritionService ports.NutritionService) *NutritionHandler {
	return &NutritionHandler{nutritionService: nutritionService}
}

func (h *NutritionHandler) GetTodayNutrition(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	nutrition, err := h.nutritionService.GetTodayNutrition(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener nutrición de hoy: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "resumen nutricional obtenido",
		"data":    nutrition,
	})
}

func (h *NutritionHandler) AddWater(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	var input ports.AddWaterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cantidad de agua requerida (+150, +250, +500): " + err.Error()})
		return
	}

	nutrition, err := h.nutritionService.AddWater(c.Request.Context(), userID.(string), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidNutritionInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al agregar ingesta de agua: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "ingesta de agua registrada",
		"data":    nutrition,
	})
}

func (h *NutritionHandler) AddMeal(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	var input ports.AddMealInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de comida inválidos"})
		return
	}
	meal, err := h.nutritionService.AddMeal(c.Request.Context(), userID.(string), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidMealInput) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al registrar comida: " + err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "comida registrada", "data": meal})
}
