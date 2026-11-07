package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/adapters/handlers/middlewares"
	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

type UserHandler struct {
	userService ports.UserService
}

func NewUserHandler(userService ports.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Register(c *gin.Context) {
	var input ports.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de registro inválidos: " + err.Error()})
		return
	}

	response, err := h.userService.Register(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al registrar usuario: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "usuario registrado exitosamente",
		"data":    response,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var input ports.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "datos de inicio de sesión inválidos: " + err.Error()})
		return
	}

	response, err := h.userService.Login(c.Request.Context(), input)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al iniciar sesión: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "inicio de sesión exitoso",
		"data":    response,
	})
}

func (h *UserHandler) GoogleLogin(c *gin.Context) {
	var input ports.GoogleLoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "token de google inválido o requerido: " + err.Error()})
		return
	}

	response, err := h.userService.GoogleLogin(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "error en autenticación con Google: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "autenticación con Google exitosa",
		"data":    response,
	})
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}

	user, err := h.userService.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener perfil: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "perfil obtenido exitosamente",
		"data":    user,
	})
}
