package http

import (
	"encoding/json"
	"net/http"
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/internal/modules/homepage/usecase"
	"toan-co-tra-backend/pkg/event"
)

type HomepageHandler struct {
	uc         *usecase.HomepageUsecase
	dispatcher *event.EventDispatcher
}

func NewHomepageHandler(uc *usecase.HomepageUsecase, dispatcher *event.EventDispatcher) *HomepageHandler {
	return &HomepageHandler{
		uc:         uc,
		dispatcher: dispatcher,
	}
}

// GetClasses xử lý GET /api/v1/homepage/classes
func (h *HomepageHandler) GetClasses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	classes, err := h.uc.GetClasses()
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Không thể tải danh sách lớp học")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    classes,
	})
}

// Register xử lý POST /api/v1/registrations (Đăng ký của phụ huynh)
func (h *HomepageHandler) Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	var lead domain.Lead
	if err := json.NewDecoder(r.Body).Decode(&lead); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu JSON gửi lên không hợp lệ")
		return
	}

	savedLead, err := h.uc.RegisterLead(&lead)
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Phát sự kiện đến Event Broker để đẩy Leads về màn hình Sales qua SSE
	h.dispatcher.Publish("registration.created", savedLead)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Đăng ký tư vấn học thử thành công!",
		"data":    savedLead,
	})
}

func (h *HomepageHandler) writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

// GetGalleryImages xử lý GET /api/v1/homepage/gallery
func (h *HomepageHandler) GetGalleryImages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	images, err := h.uc.GetGalleryImages()
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Không thể tải danh sách hình ảnh lớp học")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    images,
	})
}
