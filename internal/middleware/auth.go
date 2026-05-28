package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/pkg/logger"
)

// UserContextKey đại diện cho key lưu trữ thông tin session trong Context
type UserContextKey string
const SessionKey UserContextKey = "session"

// AuthMiddleware kiểm tra phiên đăng nhập và đính kèm session vào Context của request
type AuthMiddleware struct {
	sessions map[string]*domain.Session
	log      *logger.Logger
}

// NewAuthMiddleware khởi tạo AuthMiddleware mới
func NewAuthMiddleware(log *logger.Logger) *AuthMiddleware {
	m := &AuthMiddleware{
		sessions: make(map[string]*domain.Session),
		log:      log,
	}
	// Seeding mock active sessions
	m.sessions["token_teacher_123"] = &domain.Session{
		Token:     "token_teacher_123",
		Username:  "teacher1",
		Role:      domain.RoleTeacher,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	return m
}

// RequiredRole bảo vệ route, chỉ cho phép vai trò cụ thể đi qua
func (m *AuthMiddleware) RequiredRole(next http.HandlerFunc, requiredRole string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var token string

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && strings.ToLower(parts[0]) == "bearer" {
				token = parts[1]
			}
		}

		if token == "" {
			token = r.URL.Query().Get("token")
		}

		if token == "" {
			m.log.Warn("Từ chối truy cập: Thiếu Token xác thực")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"success":false,"error":"Yêu cầu xác thực phiên đăng nhập (Thiếu token)"}`))
			return
		}

		session, exists := m.sessions[token]
		if !exists || session.IsExpired() {
			m.log.Warn("Từ chối truy cập: Token đã hết hạn hoặc không tồn tại")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			_, _ = w.Write([]byte(`{"success":false,"error":"Phiên làm việc không hợp lệ hoặc đã hết hạn"}`))
			return
		}

		if requiredRole != "" && session.Role != requiredRole {
			m.log.Warn("Từ chối truy cập: Quyền hạn Role không phù hợp")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"success":false,"error":"Tài khoản không có đủ quyền truy cập tính năng này"}`))
			return
		}

		// Đính kèm Session vào Request Context an sau
		ctx := context.WithValue(r.Context(), SessionKey, session)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

