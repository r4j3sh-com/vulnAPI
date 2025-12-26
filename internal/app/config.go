package app

import (
	"os"
	"strconv"
)

type Config struct {
	Port        string
	Env         string
	Debug       bool
	TokenSecret string
	DBPath      string
}

func LoadConfig() Config {
	cfg := Config{
		Port:        getEnv("PORT", "9000"),
		Env:         getEnv("APP_ENV", "dev"),
		Debug:       getEnvBool("DEBUG", true),
		TokenSecret: getEnv("TOKEN_SECRET", "dev-secret"),
		DBPath:      getEnv("DB_PATH", "data/vulnapi.db"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvBool(key string, fallback bool) bool {
	if v := os.Getenv(key); v != "" {
		parsed, err := strconv.ParseBool(v)
		if err == nil {
			return parsed
		}
	}
	return fallback
}
