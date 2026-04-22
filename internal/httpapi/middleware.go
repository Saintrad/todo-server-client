package httpapi

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/Saintrad/todo-server-client/internal/auth"
	"github.com/Saintrad/todo-server-client/internal/richerror"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const ContextKeyUserID = "userID"

func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		c.Next() // process request

		duration := time.Since(start)

		slog.Info("HTTP request",
			"method", c.Request.Method,
			"path", c.FullPath(),
			"status", c.Writer.Status(),
			"duration_ms", duration.Milliseconds(),
		)
	}
}

type AuthMiddleware struct {
	JWT *auth.JWTManager
}

func NewAuthMiddleware(jwtManager *auth.JWTManager) *AuthMiddleware {
	return &AuthMiddleware{JWT: jwtManager}
}

func (m *AuthMiddleware) Handle() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "missing Authorization header",
				"code":  richerror.ErrCodeUnauthorized,
			})
			return
		}

		// Must be: "Bearer <token>"
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid authorization scheme",
				"code":  richerror.ErrCodeUnauthorized,
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")

		// Parse + Verify
		token, err := m.JWT.VerifyToken(tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(StatusCodeFor(err.Code), gin.H{
				"error": err.Message,
				"code":  err.Code,
			})
			return
		}

		// Extract claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "invalid token claims",
				"code":  richerror.ErrCodeUnauthorized,
			})
			return
		}

		// Extract uid
		uidFloat, ok := claims["uid"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "uid claim missing",
				"code":  richerror.ErrCodeUnauthorized,
			})
			return
		}

		userID := int(uidFloat)

		// Store in context
		c.Set(ContextKeyUserID, userID)

		c.Next()
	}
}
