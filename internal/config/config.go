package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	App         AppConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	SMTP        SMTPConfig
	JWT         JWTConfig
	GoogleOAuth GoogleOAuthConfig
}

type GoogleOAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type AppConfig struct {
	Port       string
	WAEndpoint string
}

type DatabaseConfig struct {
	Host string
	Port string
	User string
	Pass string
	Name string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	FromName string
	FromEmail string
}

type JWTConfig struct {
	Secret string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		App: AppConfig{
			Port:       getEnv("APP_PORT", "3000"),
			WAEndpoint: getEnv("WA_ENDPOINT", "http://localhost:8080/whatsapp"),
		},
		Database: DatabaseConfig{
			Host: getEnv("DB_HOST", "localhost"),
			Port: getEnv("DB_PORT", "3306"),
			User: getEnv("DB_USER", "root"),
			Pass: getEnv("DB_PASS", ""),
			Name: getEnv("DB_NAME", "go_trial"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "localhost"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
		},
		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", "smtp.gmail.com"),
			Port:     getEnvInt("SMTP_PORT", 587),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			FromName: getEnv("SMTP_FROM_NAME", "Go Trial"),
			FromEmail: getEnv("SMTP_FROM_EMAIL", ""),
		},
		JWT: JWTConfig{
			Secret: getEnv("JWT_SECRET", "super-secret-key-change-me"),
		},
		GoogleOAuth: GoogleOAuthConfig{
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", ""),
		},
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val, ok := os.LookupEnv(key); ok {
		if intVal, err := strconv.Atoi(val); err == nil {
			return intVal
		}
	}
	return fallback
}
