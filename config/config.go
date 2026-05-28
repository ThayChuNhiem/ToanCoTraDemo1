package config

import (
	"os"
)

// Config đại diện cho cấu hình hệ thống toàn cục
type Config struct {
	Port         string
	DatabaseURL  string
	JWTSecret    string
	ChatbotSecret string
}

// LoadConfig đọc các cấu hình từ biến môi trường hoặc gán giá trị mặc định tối ưu
func LoadConfig() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		// Mặc định chạy in-memory fallback nếu không có DB thật được cung cấp
		dbURL = "postgres://postgres:postgres@localhost:5432/toan_co_tra?sslmode=disable"
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "toan_co_tra_super_secret_key_2026"
	}

	chatbotSecret := os.Getenv("CHATBOT_SECRET")
	if chatbotSecret == "" {
		chatbotSecret = "zalo_fb_webhook_verification_token_2026"
	}

	return &Config{
		Port:         port,
		DatabaseURL:  dbURL,
		JWTSecret:    jwtSecret,
		ChatbotSecret: chatbotSecret,
	}
}
