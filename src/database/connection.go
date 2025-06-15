package database

import (
	"Eco_GO/src/models"
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error

	// Configurar DSN desde variables de entorno
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		getEnv("DB_USER", "root"),
		getEnv("DB_PASSWORD", ""),
		getEnv("DB_HOST", "localhost:3306"),
		getEnv("DB_NAME", "ecosurprise"),
	)

	// Configurar logger según entorno
	var logLevel logger.LogLevel = logger.Warn
	if getEnv("ENV", "development") == "development" {
		logLevel = logger.Info
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatal("Error al conectar con la base de datos:", err)
	}

	log.Println("Conexión a base de datos establecida")

	// Ejecutar migraciones paso a paso para evitar conflictos
	if err := runMigrations(); err != nil {
		log.Fatal("Error en auto-migración:", err)
	}

	log.Println("Base de datos migrada exitosamente")
}

func runMigrations() error {
	// Migrar en orden específico para evitar conflictos de claves foráneas

	// 1. Tablas base sin relaciones
	if err := DB.AutoMigrate(&models.Usuario{}); err != nil {
		return fmt.Errorf("error migrando Usuario: %v", err)
	}

	if err := DB.AutoMigrate(&models.Direccion{}); err != nil {
		return fmt.Errorf("error migrando Direccion: %v", err)
	}

	if err := DB.AutoMigrate(&models.Fecha{}); err != nil {
		return fmt.Errorf("error migrando Fecha: %v", err)
	}

	// 2. Tablas que extienden Usuario
	if err := DB.AutoMigrate(&models.Consumidor{}); err != nil {
		return fmt.Errorf("error migrando Consumidor: %v", err)
	}

	if err := DB.AutoMigrate(&models.Comerciante{}); err != nil {
		return fmt.Errorf("error migrando Comerciante: %v", err)
	}

	if err := DB.AutoMigrate(&models.Repartidor{}); err != nil {
		return fmt.Errorf("error migrando Repartidor: %v", err)
	}

	// 3. Tablas con relaciones simples
	if err := DB.AutoMigrate(&models.Telefono{}); err != nil {
		return fmt.Errorf("error migrando Telefono: %v", err)
	}

	if err := DB.AutoMigrate(&models.UsuarioDireccion{}); err != nil {
		return fmt.Errorf("error migrando UsuarioDireccion: %v", err)
	}

	if err := DB.AutoMigrate(&models.Producto{}); err != nil {
		return fmt.Errorf("error migrando Producto: %v", err)
	}

	if err := DB.AutoMigrate(&models.Puntuacion{}); err != nil {
		return fmt.Errorf("error migrando Puntuacion: %v", err)
	}

	// 4. Tablas con relaciones complejas
	if err := DB.AutoMigrate(&models.Orden{}); err != nil {
		return fmt.Errorf("error migrando Orden: %v", err)
	}

	if err := DB.AutoMigrate(&models.OrdenProducto{}); err != nil {
		return fmt.Errorf("error migrando OrdenProducto: %v", err)
	}

	if err := DB.AutoMigrate(&models.Pago{}); err != nil {
		return fmt.Errorf("error migrando Pago: %v", err)
	}

	// 5. Tablas de entrega
	if err := DB.AutoMigrate(&models.Entrega{}); err != nil {
		return fmt.Errorf("error migrando Entrega: %v", err)
	}

	if err := DB.AutoMigrate(&models.EntregaDireccion{}); err != nil {
		return fmt.Errorf("error migrando EntregaDireccion: %v", err)
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
