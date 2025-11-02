package config

import (
	"os"
)

type Config struct {
	BotToken         string
	ServerPort       string
	RussianPostLogin string
	RussianPostPass  string
}

func Load() *Config {
	return &Config{
		BotToken:         getEnv("BOT_TOKEN", ""),
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		RussianPostLogin: getEnv("RUSSIAN_POST_LOGIN", ""),
		RussianPostPass:  getEnv("RUSSIAN_POST_PASS", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
