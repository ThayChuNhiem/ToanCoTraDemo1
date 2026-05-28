package http

import (
	"encoding/json"
	"net/http"
	"toan-co-tra-backend/internal/modules/honor/usecase"
)

type HonorHandler struct {
	uc *usecase.HonorUsecase
}

func NewHonorHandler(uc *usecase.HonorUsecase) *HonorHandler {
	return &HonorHandler{uc: uc}
}

// GetHonoredStudents xử lý GET /api/v1/honor/students
func (h *HonorHandler) GetHonoredStudents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(`{"success":false,"error":"Phương thức HTTP không được hỗ trợ"}`))
		return
	}

	students, err := h.uc.GetHonoredStudents()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"success":false,"error":"Không thể lấy danh sách học sinh tuyên dương"}`))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    students,
	})
}
