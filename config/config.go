package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBUser             string
	DBPassword         string
	DBName             string
	DBHost             string
	DBPort             string
	JWTSecret          string
	AccessTokenTTLMin  int
	RefreshTokenTTLDays int
}


func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	accessTTL, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL_MINUTES"))
	if err != nil {
		accessTTL = 15 // по умолчанию
	}

	refreshTTL, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TTL_DAYS"))
	if err != nil {
		refreshTTL = 7 // по умолчанию
	}

	return &Config{
		DBUser:              os.Getenv("POSTGRES_USER"),
		DBPassword:          os.Getenv("POSTGRES_PASSWORD"),
		DBName:              os.Getenv("POSTGRES_DB"),
		DBHost:              os.Getenv("POSTGRES_HOST"),
		DBPort:              os.Getenv("POSTGRES_PORT"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		AccessTokenTTLMin:   accessTTL,
		RefreshTokenTTLDays: refreshTTL,
	}
}
