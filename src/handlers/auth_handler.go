package handlers

import (
	"Eco_GO/src/dto"
	"Eco_GO/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

func (ctrl *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := ctrl.authService.Login(req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (ctrl *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	usuario, err := ctrl.authService.Register(req)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, usuario)
}

func (ctrl *AuthHandler) ValidateToken(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusOK, gin.H{"valid": false})
		return
	}

	tokenString := authHeader[7:] // Remover "Bearer "
	_, err := ctrl.authService.ValidateToken(tokenString)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{"valid": false})
		return
	}

	// Generar nuevo token
	// Aquí podrías implementar la renovación del token si es necesario

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
		"token": tokenString, // O un nuevo token renovado
	})
}
