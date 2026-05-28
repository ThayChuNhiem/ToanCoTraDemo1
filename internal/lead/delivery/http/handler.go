package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/internal/lead/usecase"
	"toan-co-tra-backend/pkg/event"
)

// LeadHandler xử lý các API endpoint liên quan đến Quản lý Lead học sinh
type LeadHandler struct {
	uc          *usecase.LeadUsecase
	authService *AuthServiceHelper
	dispatcher  *event.EventDispatcher
}

// AuthServiceHelper là struct hỗ trợ xác thực vai trò cục bộ
type AuthServiceHelper struct {
	sessions map[string]*domain.Session
}

// NewAuthServiceHelper khởi tạo công cụ xác thực đơn giản
func NewAuthServiceHelper() *AuthServiceHelper {
	helper := &AuthServiceHelper{
		sessions: make(map[string]*domain.Session),
	}
	// Seeding mock active sessions for direct teacher & parent login ease
	helper.sessions["token_teacher_123"] = &domain.Session{
		Token:     "token_teacher_123",
		Username:  "teacher1",
		Role:      domain.RoleTeacher,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	helper.sessions["token_parent_123"] = &domain.Session{
		Token:     "token_parent_123",
		Username:  "parent1",
		Role:      domain.RoleParent,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	return helper
}

// NewLeadHandler khởi tạo Handler mới
func NewLeadHandler(uc *usecase.LeadUsecase, dispatcher *event.EventDispatcher) *LeadHandler {
	return &LeadHandler{
		uc:          uc,
		authService: NewAuthServiceHelper(),
		dispatcher:  dispatcher,
	}
}

// authenticate kiểm tra phiên đăng nhập và phân quyền an toàn
func (h *LeadHandler) authenticate(w http.ResponseWriter, r *http.Request, requiredRole string) *domain.Session {
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
		h.writeJSONError(w, http.StatusUnauthorized, "Yêu cầu xác thực phiên đăng nhập (Thiếu token)")
		return nil
	}

	session, exists := h.authService.sessions[token]
	if !exists || session.IsExpired() {
		h.writeJSONError(w, http.StatusUnauthorized, "Phiên làm việc không hợp lệ hoặc đã hết hạn")
		return nil
	}

	if requiredRole != "" && session.Role != requiredRole {
		h.writeJSONError(w, http.StatusForbidden, "Tài khoản của bạn không có đủ quyền truy cập tính năng này")
		return nil
	}

	return session
}

// Register xử lý POST /api/v1/registrations (Phụ huynh đăng ký)
func (h *LeadHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ, vui lòng sử dụng POST")
		return
	}

	var lead domain.Lead
	if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu JSON gửi lên không hợp lệ hoặc sai định dạng")
		return
	}

	savedLead, err := h.uc.RegisterNewLead(&lead)
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Đăng ký tư vấn học thử thành công!",
		"data":    savedLead,
	})
}

// GetLeads xử lý GET /api/v1/registrations (Sales lấy toàn bộ Lead)
func (h *LeadHandler) GetLeads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ, vui lòng sử dụng GET")
		return
	}

	session := h.authenticate(w, r, domain.RoleTeacher)
	if session == nil {
		return
	}

	list, err := h.uc.GetAllLeads()
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Lỗi hệ thống khi lấy danh sách đăng ký")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    list,
	})
}

// StreamLeads xử lý GET /api/v1/registrations/stream (Server-Sent Events)
func (h *LeadHandler) StreamLeads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Giao thức truyền dữ liệu thời gian thực không được hỗ trợ", http.StatusInternalServerError)
		return
	}

	session := h.authenticate(w, r, domain.RoleTeacher)
	if session == nil {
		return
	}

	clientChan := make(chan interface{})
	h.dispatcher.Subscribe("registration.created", clientChan)

	defer func() {
		h.dispatcher.Unsubscribe("registration.created", clientChan)
		close(clientChan)
	}()

	_, _ = w.Write([]byte("event: connected\ndata: {\"status\":\"ready\",\"user\":\"" + session.Username + "\"}\n\n"))
	flusher.Flush()

	for {
		select {
		case data, ok := <-clientChan:
			if !ok {
				return
			}

			lead, isLead := data.(*domain.Lead)
			if !isLead {
				continue
			}

			jsonData, err := json.Marshal(lead)
			if err != nil {
				continue
			}

			_, err = w.Write([]byte("event: registration_created\ndata: " + string(jsonData) + "\n\n"))
			if err != nil {
				return
			}
			flusher.Flush()

		case <-r.Context().Done():
			return
		}
	}
}

// ToggleConsulted xử lý POST /api/v1/registrations/toggle-consulted (Đánh dấu đã tư vấn / chưa tư vấn)
func (h *LeadHandler) ToggleConsulted(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ, vui lòng sử dụng POST")
		return
	}

	session := h.authenticate(w, r, domain.RoleTeacher)
	if session == nil {
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu JSON gửi lên không hợp lệ")
		return
	}

	list, err := h.uc.GetAllLeads()
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Lỗi lấy danh sách đăng ký")
		return
	}

	var found *domain.Lead
	for _, lead := range list {
		if lead.ID == req.ID {
			found = lead
			break
		}
	}

	if found == nil {
		h.writeJSONError(w, http.StatusNotFound, "Không tìm thấy thông tin đăng ký tương ứng")
		return
	}

	newStatus := "contacted"
	if found.Status == "contacted" {
		newStatus = "pending"
	}

	err = h.uc.UpdateLeadStatus(req.ID, newStatus)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Lỗi hệ thống khi cập nhật trạng thái")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"status":  newStatus,
	})
}

// Login xử lý POST /api/v1/auth/login
func (h *LeadHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ, vui lòng sử dụng POST")
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu JSON đăng nhập không hợp lệ")
		return
	}

	var session *domain.Session
	var name string

	if creds.Username == "teacher1" && creds.Password == "password123" {
		session = h.authService.sessions["token_teacher_123"]
		name = "Cô Giáo Trà"
	} else if creds.Username == "parent1" && creds.Password == "password123" {
		session = h.authService.sessions["token_parent_123"]
		name = "Chị Nguyễn Lan (Mẹ bé Đức)"
	} else {
		h.writeJSONError(w, http.StatusUnauthorized, "Tên đăng nhập hoặc mật khẩu không chính xác")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Đăng nhập thành công!",
		"data": map[string]interface{}{
			"token":    session.Token,
			"username": session.Username,
			"name":     name,
			"role":     session.Role,
		},
	})
}

// Logout xử lý POST /api/v1/auth/logout
func (h *LeadHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ, vui lòng sử dụng POST")
		return
	}

	token := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = token[7:]
	}

	if token == "" {
		h.writeJSONError(w, http.StatusBadRequest, "Thiếu token đăng xuất")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Đăng xuất thành công!",
	})
}

// writeJSONError helper to return consistent JSON error models
func (h *LeadHandler) writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

