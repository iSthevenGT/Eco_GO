package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Config struct {
	Port           string
	MYSQL_URL      string
	MYSQL_DB       string
	MYSQL_USER     string
	MYSQL_PASSWORD string
	MYSQL_HOST     string
	MYSQL_PORT     string
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file:")
	}

	return Config{
		Port:           getEnv("PORT", "8080"),
		MYSQL_URL:      getEnv("MYSQL_URL", ""),
		MYSQL_DB:       getEnv("MYSQL_DB", ""),
		MYSQL_USER:     getEnv("MYSQL_USER", ""),
		MYSQL_PASSWORD: getEnv("MYSQL_PASSWORD", ""),
		MYSQL_HOST:     getEnv("MYSQL_HOST", "localhost"),
		MYSQL_PORT:     getEnv("MYSQL_PORT", "3306"),
	}
}

func ConectToDB(cfg Config) *gorm.DB {

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.MYSQL_USER,
		cfg.MYSQL_PASSWORD,
		cfg.MYSQL_HOST,
		cfg.MYSQL_PORT,
		cfg.MYSQL_DB,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Cannot connect to DB:", err)
	}
	return db
}
