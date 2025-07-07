package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBURL               string
	JWTSecret           string
	AccessTokenTTLMin   int
	RefreshTokenTTLDays int
}

func LoadConfig() *Config {
	accessTTL, err := strconv.Atoi(os.Getenv("ACCESS_TOKEN_TTL_MINUTES"))
	if err != nil {
		accessTTL = 15
	}

	refreshTTL, err := strconv.Atoi(os.Getenv("REFRESH_TOKEN_TTL_DAYS"))
	if err != nil {
		refreshTTL = 7
	}

	return &Config{
		DBURL:               os.Getenv("DB_URL"),
		JWTSecret:           os.Getenv("JWT_SECRET"),
		AccessTokenTTLMin:   accessTTL,
		RefreshTokenTTLDays: refreshTTL,
	}
}
