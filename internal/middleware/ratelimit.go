package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// IPClient đại diện cho thông tin bucket giới hạn tốc độ của từng IP Client
type IPClient struct {
	tokens     float64
	lastAccess time.Time
}

// RateLimiter ngăn chặn tấn công spam hoặc brute force từ một IP Client cụ thể
type RateLimiter struct {
	mu      sync.Mutex
	clients map[string]*IPClient
	rate    float64 // Tốc độ nạp lại: số token nạp lại mỗi giây
	max     float64 // Tối đa số token trong bucket (tương đương số request cho phép bắn dồn dập)
}

// NewRateLimiter khởi tạo bộ lọc giới hạn với tốc độ và dung lượng bucket thiết lập
func NewRateLimiter(requestsPerSecond float64, maxBurst float64) *RateLimiter {
	return &RateLimiter{
		clients: make(map[string]*IPClient),
		rate:    requestsPerSecond,
		max:     maxBurst,
	}
}

// LimitMiddleware là hàm chặn trung gian HTTP để kiểm tra IP
func (rl *RateLimiter) LimitMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 1. Phân tích địa chỉ IP của Client gửi yêu cầu lên
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			ip = r.RemoteAddr
		}

		// 2. Kiểm tra token khả dụng
		if !rl.allow(ip) {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusTooManyRequests)
			_, _ = w.Write([]byte(`{"success":false,"error":"Hệ thống phát hiện tần suất gửi yêu cầu quá nhanh từ IP của bạn. Vui lòng thử lại sau ít phút!"}`))
			return
		}

		next.ServeHTTP(w, r)
	}
}

// allow triển khai thuật toán Token Bucket luồng an toàn (Thread-safe)
func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	client, exists := rl.clients[ip]
	if !exists {
		rl.clients[ip] = &IPClient{
			tokens:     rl.max,
			lastAccess: now,
		}
		return true
	}

	// Tính toán lượng token nạp lại dựa trên thời gian trôi qua
	elapsed := now.Sub(client.lastAccess).Seconds()
	client.lastAccess = now

	client.tokens += elapsed * rl.rate
	if client.tokens > rl.max {
		client.tokens = rl.max
	}

	// Kiểm tra nếu bucket còn token để tiêu thụ
	if client.tokens >= 1.0 {
		client.tokens -= 1.0
		return true
	}

	return false
}
