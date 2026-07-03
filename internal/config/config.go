package config

import (
	"os"
	"strconv"
)

type Config struct {
    Port     string
    DBHost   string
    DBPort   int
    DBUser   string
    DBPass   string
    DBName   string
	JWTSecret   string
}

func Load() Config {

	portStr := os.Getenv("DB_PORT")

	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 5432
	}

    return Config{
        Port:   getEnv("PORT", "8080"),
        DBHost: getEnv("DB_HOST", "localhost"),
        DBPort: port,
        DBUser: getEnv("DB_USER", "user"),
        DBPass: getEnv("DB_PASS", "password"),
        DBName: getEnv("DB_NAME", "todo"),
		JWTSecret: os.Getenv("JWT_SECRET"),
    }
}

func getEnv(key, fallback string) string {
    if v := os.Getenv(key); v != "" {
        return v
    }
    return fallback
}