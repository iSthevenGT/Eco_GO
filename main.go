package main

import (
	"Eco_GO/src/config"
	"Eco_GO/src/database"
	"Eco_GO/src/routes"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Configurar Gin
	if os.Getenv("ENV") == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar base de datos
	database.InitDB()

	// Configurar router
	router := gin.Default()

	// Configurar CORS
	config.SetupCORS(router)

	// Configurar rutas
	routes.SetupRoutes(router)

	// Iniciar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Servidor iniciando en puerto %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Error al iniciar servidor:", err)
	}
}
