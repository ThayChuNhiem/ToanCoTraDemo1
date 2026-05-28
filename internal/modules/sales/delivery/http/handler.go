package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/internal/modules/sales/usecase"
	"toan-co-tra-backend/internal/repository"
	"toan-co-tra-backend/pkg/event"
)

type SalesHandler struct {
	uc         *usecase.SalesUsecase
	repo       *repository.CentralRepository
	dispatcher *event.EventDispatcher
	sessions   map[string]*domain.Session
	mu         sync.RWMutex
}

func NewSalesHandler(uc *usecase.SalesUsecase, repo *repository.CentralRepository, dispatcher *event.EventDispatcher) *SalesHandler {
	h := &SalesHandler{
		uc:         uc,
		repo:       repo,
		dispatcher: dispatcher,
		sessions:   make(map[string]*domain.Session),
	}

	// Seeding mock active sessions for ease of initial API request bypassing
	h.sessions["token_teacher_123"] = &domain.Session{
		Token:     "token_teacher_123",
		Username:  "teacher1",
		Role:      domain.RoleTeacher,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}
	h.sessions["token_admin_123"] = &domain.Session{
		Token:     "token_admin_123",
		Username:  "admin1",
		Role:      domain.RoleAdmin,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	return h
}

// GetLeads xử lý GET /api/v1/registrations
func (h *SalesHandler) GetLeads(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	// Kiểm tra phân quyền: Chỉ Teacher hoặc Admin mới được xem danh sách Leads
	session := h.authenticate(w, r)
	if session == nil {
		return
	}

	if session.Role != domain.RoleTeacher && session.Role != domain.RoleAdmin {
		h.writeJSONError(w, http.StatusForbidden, "Tài khoản không có quyền truy cập thông tin")
		return
	}

	leads, err := h.uc.GetAllLeads()
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Lỗi khi lấy danh sách đăng ký")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    leads,
	})
}

// ToggleConsulted xử lý POST /api/v1/registrations/toggle-consulted (Đánh dấu đã tư vấn / chưa tư vấn)
func (h *SalesHandler) ToggleConsulted(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	session := h.authenticate(w, r)
	if session == nil {
		return
	}

	if session.Role != domain.RoleTeacher && session.Role != domain.RoleAdmin {
		h.writeJSONError(w, http.StatusForbidden, "Tài khoản không có quyền thực hiện thao tác này")
		return
	}

	var req struct {
		ID int64 `json:"id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ")
		return
	}

	leads, err := h.uc.GetAllLeads()
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Lỗi lấy danh sách")
		return
	}

	var found *domain.Lead
	for _, l := range leads {
		if l.ID == req.ID {
			found = l
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

	err = h.uc.UpdateLeadStatus(req.ID, newStatus, session.Username)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Lỗi khi cập nhật trạng thái tư vấn")
		return
	}

	// Đẩy trạng thái cập nhật qua SSE tới các màn hình sales khác đang xem
	found.Status = newStatus
	found.ConsultedBy = session.Username
	if newStatus == "pending" {
		found.ConsultedBy = ""
	}
	h.dispatcher.Publish("registration.updated", found)


	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"status":  newStatus,
	})
}

// StreamLeads xử lý Server-Sent Events (SSE) thời gian thực
func (h *SalesHandler) StreamLeads(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Giao thức truyền dữ liệu thời gian thực không được hỗ trợ", http.StatusInternalServerError)
		return
	}

	session := h.authenticate(w, r)
	if session == nil {
		return
	}

	clientChan := make(chan interface{})
	h.dispatcher.Subscribe("registration.created", clientChan)
	h.dispatcher.Subscribe("registration.updated", clientChan)

	defer func() {
		h.dispatcher.Unsubscribe("registration.created", clientChan)
		h.dispatcher.Unsubscribe("registration.updated", clientChan)
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

// Login xử lý POST /api/v1/auth/login
func (h *SalesHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu đăng nhập không hợp lệ")
		return
	}

	// Tra cứu tài khoản trực tiếp trong cơ sở dữ liệu central_repo
	user, err := h.repo.FindUser(creds.Username)
	if err != nil || user.Password != creds.Password {
		h.writeJSONError(w, http.StatusUnauthorized, "Tên đăng nhập hoặc mật khẩu không chính xác")
		return
	}

	// Tạo Token ngẫu nhiên hoặc gán cố định để đơn giản hóa
	token := fmt.Sprintf("token_%s_%d", user.Username, time.Now().UnixNano())
	if user.Username == "teacher1" {
		token = "token_teacher_123"
	} else if user.Username == "admin1" {
		token = "token_admin_123"
	}

	session := &domain.Session{
		Token:     token,
		Username:  user.Username,
		Role:      user.Role,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	}

	h.mu.Lock()
	h.sessions[token] = session
	h.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Đăng nhập thành công!",
		"data": map[string]interface{}{
			"token":    session.Token,
			"username": session.Username,
			"name":     user.Name,
			"role":     session.Role,
		},
	})
}

// Logout xử lý POST /api/v1/auth/logout
func (h *SalesHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	token := r.Header.Get("Authorization")
	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = token[7:]
	}

	if token == "" {
		token = r.URL.Query().Get("token")
	}

	if token != "" {
		h.mu.Lock()
		delete(h.sessions, token)
		h.mu.Unlock()
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Đăng xuất thành công!",
	})
}

// authenticate hỗ trợ lấy Session hiện tại qua Bearer Token
func (h *SalesHandler) authenticate(w http.ResponseWriter, r *http.Request) *domain.Session {
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

	h.mu.RLock()
	session, exists := h.sessions[token]
	h.mu.RUnlock()

	if !exists || session.IsExpired() {
		h.writeJSONError(w, http.StatusUnauthorized, "Phiên làm việc không hợp lệ hoặc đã hết hạn")
		return nil
	}

	return session
}

func (h *SalesHandler) GetSession(token string) *domain.Session {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.sessions[token]
}

func (h *SalesHandler) AddSession(token string, session *domain.Session) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sessions[token] = session
}

func (h *SalesHandler) writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}
