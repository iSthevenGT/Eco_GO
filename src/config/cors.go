package config

import (
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupCORS(router *gin.Engine) {
	config := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	// Configurar orígenes permitidos
	ip := os.Getenv("APP_IP")
	if ip == "*" {
		config.AllowAllOrigins = true
	} else {
		allowedOrigins := []string{
			"http://localhost:8081",
			"http://localhost:3000",
		}
		if ip != "" {
			allowedOrigins = append(allowedOrigins,
				"http://"+ip+":8081",
				"exp://"+ip+":8081",
			)
		}
		config.AllowOrigins = allowedOrigins
	}

	router.Use(cors.New(config))
}
