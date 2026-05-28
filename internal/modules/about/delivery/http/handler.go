package http

import (
	"encoding/json"
	"net/http"
	"toan-co-tra-backend/internal/modules/about/usecase"
)

type AboutHandler struct {
	uc *usecase.AboutUsecase
}

func NewAboutHandler(uc *usecase.AboutUsecase) *AboutHandler {
	return &AboutHandler{uc: uc}
}

// GetTeachers xử lý GET /api/v1/about/teachers
func (h *AboutHandler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"success":false,"error":"Phương thức HTTP không được hỗ trợ"}`))
		return
	}

	teachers, err := h.uc.GetTeachers()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"success":false,"error":"Không thể lấy danh sách giáo viên"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    teachers,
	})
}
