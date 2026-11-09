package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port               string
	Env                string
	DBHost             string
	DBPort             string
	DBUser             string
	DBPassword         string
	DBName             string
	JWTSecret          string
	JWTExpirationHours int
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found; reading configuration from environment.")
	}

	jwtExpHours, err := strconv.Atoi(getEnv("JWT_EXPIRATION_HOURS", "12"))
	if err != nil || jwtExpHours < 1 {
		jwtExpHours = 24
	}

	return &Config{
		Port:               getEnv("PORT", "8080"),
		Env:                getEnv("ENV", "development"),
		DBHost:             getEnv("DB_HOST", "127.0.0.1"),
		DBPort:             getEnv("DB_PORT", "3306"),
		DBUser:             getEnv("DB_USER", "root"),
		DBPassword:         getEnv("DB_PASSWORD", ""),
		DBName:             getEnv("DB_NAME", "forgefit_db"),
		JWTSecret:          getEnv("JWT_SECRET", "default-jwt-secret-key-change-me"),
		JWTExpirationHours: jwtExpHours,
	}
}

func (c *Config) Validate() error {
	if c.Port == "" {
		return fmt.Errorf("PORT must not be empty")
	}
	if c.Env != "production" {
		return nil
	}
	if c.DBHost == "" || c.DBPort == "" || c.DBUser == "" || c.DBPassword == "" || c.DBName == "" {
		return fmt.Errorf("DB_HOST, DB_PORT, DB_USER, DB_PASSWORD and DB_NAME are required in production")
	}
	if len(c.JWTSecret) < 32 || c.JWTSecret == "default-jwt-secret-key-change-me" {
		return fmt.Errorf("JWT_SECRET must be a unique secret of at least 32 characters in production")
	}
	return nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
