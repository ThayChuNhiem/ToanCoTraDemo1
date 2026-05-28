package chatbot

import (
	"encoding/json"
	"fmt"
	"net/http"
	"toan-co-tra-backend/config"
	"toan-co-tra-backend/internal/lead/usecase"
	"toan-co-tra-backend/pkg/logger"
)

// WebhookHandler xử lý các dữ liệu webhook đẩy về từ API Zalo/Facebook
type WebhookHandler struct {
	uc     *usecase.LeadUsecase
	cfg    *config.Config
	log    *logger.Logger
}

// NewWebhookHandler khởi tạo WebhookHandler mới
func NewWebhookHandler(uc *usecase.LeadUsecase, cfg *config.Config, log *logger.Logger) *WebhookHandler {
	return &WebhookHandler{
		uc:     uc,
		cfg:    cfg,
		log:    log,
	}
}

// VerifyWebhook xử lý yêu cầu xác thực webhook ban đầu từ Zalo/Facebook (ví dụ: GET challenge token)
func (h *WebhookHandler) VerifyWebhook(w http.ResponseWriter, r *http.Request) {
	// Lấy token xác thực gửi từ platform
	verifyToken := r.URL.Query().Get("hub.verify_token")
	challenge := r.URL.Query().Get("hub.challenge")

	if verifyToken == "" {
		h.log.Warn("Yêu cầu xác thực Webhook thiếu verify_token")
		http.Error(w, "Thiếu token xác thực", http.StatusBadRequest)
		return
	}

	// So khớp với token đã được thiết lập bảo mật trong config
	if verifyToken != h.cfg.ChatbotSecret {
		h.log.Error("Xác thực Webhook thất bại: Token không khớp", fmt.Errorf("nhận được: %s", verifyToken))
		http.Error(w, "Xác thực token thất bại", http.StatusForbidden)
		return
	}

	h.log.Info("Xác thực Webhook kết nối Zalo/Facebook thành công!")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(challenge))
}

// HandleIncomingMessages xử lý gói tin JSON chứa tin nhắn người dùng đẩy về qua POST
func (h *WebhookHandler) HandleIncomingMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Phương thức HTTP không hỗ trợ", http.StatusMethodNotAllowed)
		return
	}

	// Phân tích dữ liệu sự kiện gửi từ Zalo/Facebook
	var event struct {
		SenderID  string `json:"sender_id"`
		Message   string `json:"message"`
		Timestamp int64  `json:"timestamp"`
	}

	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		h.log.Error("Lỗi phân tích JSON webhook tin nhắn", err)
		http.Error(w, "Dữ liệu JSON không hợp lệ", http.StatusBadRequest)
		return
	}

	h.log.Info(fmt.Sprintf("Nhận tin nhắn chatbot từ Sender [%s]: %s", event.SenderID, event.Message))

	// Trả về kết quả xác nhận đã nhận sự kiện thành công (Zalo/Facebook yêu cầu phản hồi 200 OK nhanh để tránh timeout)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Đã nhận gói tin webhook thành công",
	})
}
