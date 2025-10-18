package infrastructure

import "os"

type Config struct {
	ServerPort string
	RedisAddr  string
}

func LoadConfig() Config {
	return Config{
		ServerPort: getEnv("SERVER_PORT", "8080"),
		RedisAddr:  getEnv("REDIS_ADDR", "localhost:6379"),
	}
}

func getEnv(key, defaultVal string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultVal
}
