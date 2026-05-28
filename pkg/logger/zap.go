package logger

import (
	"fmt"
	"time"
)

// Logger đại diện cho trình ghi nhật ký đơn giản nhưng định dạng đẹp, giúp quản lý thông tin ghi log
type Logger struct {
	ServiceName string
}

// NewZapLogger khởi tạo Logger hệ thống
func NewZapLogger(serviceName string) *Logger {
	return &Logger{ServiceName: serviceName}
}

// Info ghi nhận thông tin hoạt động thường nhật
func (l *Logger) Info(message string) {
	fmt.Printf(" [INFO]  %s | %s | %s\n", time.Now().Format("2006-01-02 15:04:05"), l.ServiceName, message)
}

// Warn ghi nhận cảnh báo nguy cơ hệ thống
func (l *Logger) Warn(message string) {
	fmt.Printf("⚠️ [WARN]  %s | %s | %s\n", time.Now().Format("2006-01-02 15:04:05"), l.ServiceName, message)
}

// Error ghi nhận các lỗi phát sinh nghiêm trọng
func (l *Logger) Error(message string, err error) {
	errStr := "nil"
	if err != nil {
		errStr = err.Error()
	}
	fmt.Printf("❌ [ERROR] %s | %s | %s | Error Details: %s\n", time.Now().Format("2006-01-02 15:04:05"), l.ServiceName, message, errStr)
}
