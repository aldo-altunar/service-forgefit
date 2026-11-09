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

func (h *UserHandler) GetAccountSettings(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	user, err := h.userService.GetProfile(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ajustes de cuenta obtenidos", "data": user})
}

func (h *UserHandler) UpdateAccountSettings(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	var input ports.UpdateAccountSettingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ajustes de cuenta inválidos: " + err.Error()})
		return
	}
	user, err := h.userService.UpdateAccountSettings(c.Request.Context(), userID.(string), input)
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al actualizar ajustes: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "ajustes de cuenta actualizados", "data": user})
}

func (h *UserHandler) ChangePasswordDisabled(c *gin.Context) {
	c.JSON(http.StatusNotImplemented, gin.H{"error": "el cambio de contraseña está deshabilitado"})
}

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, exists := c.Get(middlewares.ContextUserIDKey)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
		return
	}
	if err := h.userService.DeleteAccount(c.Request.Context(), userID.(string)); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al eliminar cuenta: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cuenta eliminada y datos personales anonimizados"})
}
