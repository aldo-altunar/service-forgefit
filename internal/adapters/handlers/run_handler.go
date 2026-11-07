package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/adapters/handlers/middlewares"
	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type RunHandler struct {
	runService ports.RunService
}

func NewRunHandler(runService ports.RunService) *RunHandler {
	return &RunHandler{runService: runService}
}

func (h *RunHandler) StartRun(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	var input ports.StartRunInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de carrera inválidos"})
		return
	}

	run, err := h.runService.StartRun(c.Request.Context(), userID.(string), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidGPSPoint) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al iniciar carrera: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "carrera iniciada",
		"data":    run,
	})
}

func (h *RunHandler) GetActiveRun(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	run, err := h.runService.GetActiveRun(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener carrera en curso: " + err.Error()})
		return
	}

	if run == nil {
		c.JSON(http.StatusOK, gin.H{"data": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "carrera en curso obtenida",
		"data":    run,
	})
}

func (h *RunHandler) GetRunHistory(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	items, err := h.runService.GetRunHistory(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener historial: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items})
}

func (h *RunHandler) PauseRun(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	runID := c.Param("id")

	run, err := h.runService.PauseRun(c.Request.Context(), userID.(string), runID)
	if err != nil {
		if errors.Is(err, domain.ErrRunNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrInvalidRunState) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al pausar carrera: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "carrera pausada",
		"data":    run,
	})
}

func (h *RunHandler) ResumeRun(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	runID := c.Param("id")

	run, err := h.runService.ResumeRun(c.Request.Context(), userID.(string), runID)
	if err != nil {
		if errors.Is(err, domain.ErrRunNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrInvalidRunState) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al reanudar carrera: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "carrera reanudada",
		"data":    run,
	})
}

func (h *RunHandler) FinishRun(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	runID := c.Param("id")

	var input ports.UpdateRunMetricsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "métricas inválidas"})
		return
	}

	run, err := h.runService.FinishRun(c.Request.Context(), userID.(string), runID, input)
	if err != nil {
		if errors.Is(err, domain.ErrRunNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrInvalidRunState) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al finalizar carrera: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "carrera finalizada",
		"data":    run,
	})
}

func (h *RunHandler) AddGPSPoint(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	runID := c.Param("id")

	var input ports.AddGPSPointInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos GPS inválidos: " + err.Error()})
		return
	}

	point, err := h.runService.AddGPSPoint(c.Request.Context(), userID.(string), runID, input)
	if err != nil {
		if errors.Is(err, domain.ErrRunNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrInvalidGPSPoint) {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if errors.Is(err, domain.ErrInvalidRunState) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al registrar punto GPS: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "punto GPS registrado",
		"data":    point,
	})
}
