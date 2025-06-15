package middleware

import (
	"Eco_GO/src/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	authService := services.NewAuthService()

	return func(c *gin.Context) {
		// Permitir solicitudes OPTIONS
		if c.Request.Method == "OPTIONS" {
			c.Next()
			return
		}

		// Permitir rutas públicas
		if strings.HasPrefix(c.Request.URL.Path, "/api/auth/") ||
			strings.HasPrefix(c.Request.URL.Path, "/usuarios/") ||
			strings.HasPrefix(c.Request.URL.Path, "/productos/") {
			c.Next()
			return
		}

		// Extraer token del header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token de autorización requerido"})
			c.Abort()
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Validar token
		claims, err := authService.ValidateToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
			c.Abort()
			return
		}

		// Extraer información del usuario del token
		userID := uint((*claims)["id"].(float64))
		role := (*claims)["rol"].(string)

		// Guardar información en el contexto
		c.Set("userID", userID)
		c.Set("role", role)

		c.Next()
	}
}

func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Rol no encontrado"})
			c.Abort()
			return
		}

		roleString := userRole.(string)
		for _, role := range roles {
			if roleString == role {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{"error": "Permisos insuficientes"})
		c.Abort()
	}
}
