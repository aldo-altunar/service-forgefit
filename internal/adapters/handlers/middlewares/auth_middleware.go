package middlewares

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"ForgeFit_Servidor/internal/core/domain"
	"ForgeFit_Servidor/internal/core/ports"
)

// BootstrapAdminTokenMiddleware allows a one-time admin bootstrap only while
// the server has an explicitly configured secret. The repository enforces that
// the first admin does not already exist.
func BootstrapAdminTokenMiddleware(configuredToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if configuredToken == "" {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "ruta no disponible"})
			return
		}
		providedToken := c.GetHeader("X-Admin-Bootstrap-Token")
		if subtle.ConstantTimeCompare([]byte(providedToken), []byte(configuredToken)) != 1 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "token de bootstrap inválido"})
			return
		}
		c.Next()
	}
}

const (
	AuthorizationHeaderKey  = "Authorization"
	AuthorizationTypeBearer = "bearer"
	ContextUserIDKey        = "userID"
	ContextUserEmailKey     = "userEmail"
	ContextTargetUserIDKey  = "targetUserID"
)

func AuthMiddleware(tokenService ports.TokenService, userRepository ports.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AuthorizationHeaderKey)
		if len(authHeader) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "authorization header is required",
			})
			return
		}

		fields := strings.Fields(authHeader)
		if len(fields) < 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization header format",
			})
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != AuthorizationTypeBearer {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "unsupported authorization type",
			})
			return
		}

		accessToken := fields[1]
		claims, err := tokenService.ValidateToken(accessToken)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid or expired token",
			})
			return
		}
		active, err := userRepository.IsAccountActive(c.Request.Context(), claims.UserID)
		if err != nil || !active {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "account is unavailable",
			})
			return
		}

		c.Set(ContextUserIDKey, claims.UserID)
		c.Set(ContextUserEmailKey, claims.Email)
		c.Next()
	}
}

// AdminMiddleware checks the persisted role so role changes take effect without
// waiting for an existing JWT to expire.
func AdminMiddleware(userRepository ports.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get(ContextUserIDKey)
		if !exists {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "usuario no autenticado"})
			return
		}
		user, err := userRepository.FindByID(c.Request.Context(), userID.(string))
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "cuenta no disponible"})
			return
		}
		if user.Role != domain.UserRoleAdmin {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "se requiere rol de administrador"})
			return
		}
		c.Next()
	}
}

// TargetUserMiddleware validates the recipient ID on admin assignment routes.
func TargetUserMiddleware(userRepository ports.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		targetUserID := c.Param("userId")
		if targetUserID == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "usuario destino requerido"})
			return
		}
		user, err := userRepository.FindByID(c.Request.Context(), targetUserID)
		if err != nil || user == nil {
			c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "usuario destino no encontrado"})
			return
		}
		c.Set(ContextTargetUserIDKey, targetUserID)
		c.Next()
	}
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-Admin-Bootstrap-Token, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Next()
	}
}
