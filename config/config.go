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
		// Mặc định chạy CSDL Postgres đám mây Neon thực tế nếu không có biến môi trường
		dbURL = "postgresql://neondb_owner:npg_JkyIHBd7l4YX@ep-spring-field-aqpgqquh.c-8.us-east-1.aws.neon.tech/neondb?sslmode=require"
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
