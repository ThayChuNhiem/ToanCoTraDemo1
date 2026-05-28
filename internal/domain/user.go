package domain

import (
	"time"
)

// Khai báo các vai trò hệ thống
const (
	RoleAdmin   = "admin"
	RoleTeacher = "teacher"
	RoleParent  = "parent"
)

// User đại diện cho thực thể Tài khoản người dùng
type User struct {
	Username string `json:"username"`
	Password string `json:"password"` // Để serialize nhận dạng trong admin panel hoặc ẩn đi khi cần
	Name     string `json:"name"`
	Role     string `json:"role"` // "admin", "teacher" hoặc "parent"
}

// Session đại diện cho phiên làm việc hoạt động của người dùng
type Session struct {
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
}

// IsExpired kiểm tra xem Session đã hết hạn chưa
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// UserRepository định nghĩa cách thức tương tác lưu trữ cho Tài khoản hệ thống
type UserRepository interface {
	FindUser(username string) (*User, error)
	SaveUser(user *User) error
	FindAllUsers() ([]*User, error)
	DeleteUser(username string) error
}


