package infrastructure

import "os"

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
}

type Config struct {
	ServerPort string
	RedisAddr  string

	SMTP SMTPConfig
}

func LoadConfig() Config {
	return Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),

		SMTP: SMTPConfig{
			Host:     getEnv("SMTP_HOST", ""),
			Port:     getEnv("SMTP_PORT", "587"),
			Username: getEnv("SMTP_USERNAME", ""),
			Password: getEnv("SMTP_PASSWORD", ""),
			From:     getEnv("SMTP_FROM", ""),
		},
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
