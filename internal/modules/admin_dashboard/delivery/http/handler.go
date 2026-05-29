package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/internal/modules/admin_dashboard/usecase"
	salesHttp "toan-co-tra-backend/internal/modules/sales/delivery/http"
)

type AdminHandler struct {
	uc           *usecase.AdminUsecase
	salesHandler *salesHttp.SalesHandler
}

func NewAdminHandler(uc *usecase.AdminUsecase, salesHandler *salesHttp.SalesHandler) *AdminHandler {
	return &AdminHandler{
		uc:           uc,
		salesHandler: salesHandler,
	}
}

// authenticate kiểm duyệt phân quyền bắt buộc là Admin
func (h *AdminHandler) authenticateAdmin(w http.ResponseWriter, r *http.Request) *domain.Session {
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
		h.writeJSONError(w, http.StatusUnauthorized, "Yêu cầu xác thực phiên đăng nhập Admin (Thiếu token)")
		return nil
	}

	session := h.salesHandler.GetSession(token)
	if session == nil || session.IsExpired() {
		h.writeJSONError(w, http.StatusUnauthorized, "Phiên đăng nhập không hợp lệ hoặc đã hết hạn")
		return nil
	}

	if session.Role != domain.RoleAdmin {
		h.writeJSONError(w, http.StatusForbidden, "Truy cập bị từ chối: Chỉ tài khoản Admin mới có quyền truy cập tính năng này")
		return nil
	}

	return session
}

// GetReports xử lý GET /api/v1/admin/reports
func (h *AdminHandler) GetReports(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	report, err := h.uc.GetSalesReport()
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Lỗi tổng hợp dữ liệu báo cáo")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    report,
	})
}

// UpdateClassPrice xử lý POST /api/v1/admin/classes/update
func (h *AdminHandler) UpdateClassPrice(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	var req struct {
		ID    string `json:"id"`
		Price string `json:"price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu JSON gửi lên không hợp lệ")
		return
	}

	err := h.uc.UpdateClassPrice(req.ID, req.Price)
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Cập nhật học phí lớp học thành công!",
	})
}

// UpdateTeacherBio xử lý POST /api/v1/admin/teachers/update
func (h *AdminHandler) UpdateTeacherBio(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	var req domain.Teacher
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu giảng viên không hợp lệ")
		return
	}

	err := h.uc.UpdateTeacherBio(&req)
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Cập nhật tiểu sử giáo viên thành công!",
	})
}

// ManageStudent xử lý POST /api/v1/admin/students/manage (Thêm mới hoặc Xóa)
func (h *AdminHandler) ManageStudent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	var req struct {
		Action      string          `json:"action"` // "add" hoặc "delete"
		ID          int64           `json:"id"`
		StudentData *domain.Student `json:"student_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ")
		return
	}

	if req.Action == "add" {
		if req.StudentData == nil {
			h.writeJSONError(w, http.StatusBadRequest, "Thiếu dữ liệu học sinh cần thêm mới")
			return
		}

		err := h.uc.AddHonoredStudent(req.StudentData)
		if err != nil {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Thêm học sinh vinh danh bảng vàng thành công!",
			"data":    req.StudentData,
		})
		return
	} else if req.Action == "delete" {
		err := h.uc.DeleteHonoredStudent(req.ID)
		if err != nil {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Xóa học sinh vinh danh thành công!",
		})
		return
	} else {
		h.writeJSONError(w, http.StatusBadRequest, "Thao tác không được hỗ trợ (chỉ chấp nhận 'add' hoặc 'delete')")
	}
}

// CreateTeacherAccount xử lý POST /api/v1/admin/users/create (CHỈ ADMIN MỚI ĐƯỢC TẠO SALES)
func (h *AdminHandler) CreateTeacherAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	// Xác thực quyền Admin tuyệt đối
	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	var req domain.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu tài khoản không hợp lệ")
		return
	}

	err := h.uc.CreateTeacherAccount(&req)
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": fmt.Sprintf("Tạo mới tài khoản Giáo viên/Sales '%s' thành công!", req.Username),
	})
}

// ListUsers xử lý GET /api/v1/admin/users
func (h *AdminHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	users, err := h.uc.GetAllUsers()
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Không thể lấy danh sách tài khoản")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    users,
	})
}

// DeleteUser xử lý POST /api/v1/admin/users/delete (Admin xóa tài khoản sales)
func (h *AdminHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ")
		return
	}

	err := h.uc.DeleteSalesUser(req.Username)
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Xóa tài khoản thành công!",
	})
}

// ManageTeacher xử lý POST /api/v1/admin/teachers/manage (Admin CRUD đồng nghiệp giáo viên)
func (h *AdminHandler) ManageTeacher(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	var req struct {
		Action      string          `json:"action"` // "add" hoặc "delete"
		ID          int64           `json:"id"`
		TeacherData *domain.Teacher `json:"teacher_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ")
		return
	}

	if req.Action == "add" {
		if req.TeacherData == nil {
			h.writeJSONError(w, http.StatusBadRequest, "Thiếu dữ liệu giáo viên cần thêm mới")
			return
		}

		err := h.uc.AddTeacher(req.TeacherData)
		if err != nil {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Thêm giáo viên mới thành công!",
			"data":    req.TeacherData,
		})
		return
	} else if req.Action == "delete" {
		err := h.uc.DeleteTeacher(req.ID)
		if err != nil {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Xóa thông tin giáo viên thành công khỏi trang Giới Thiệu!",
		})
		return
	} else {
		h.writeJSONError(w, http.StatusBadRequest, "Thao tác không được hỗ trợ (chỉ chấp nhận 'add' hoặc 'delete')")
	}
}


// UploadFile xử lý POST /api/v1/admin/upload (Admin upload ảnh giáo viên & học sinh)
func (h *AdminHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	// Xác thực quyền Admin tuyệt đối
	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	// Parse form dữ liệu lên đến 10MB
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu tệp tin tải lên quá lớn (tối đa 10MB)")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Không tìm thấy tệp tin upload 'file'")
		return
	}
	defer file.Close()

	// Đảm bảo thư mục upload tồn tại
	uploadDir := "./assets/uploads"
	err = os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Lỗi hệ thống khi khởi tạo thư mục lưu trữ ảnh")
		return
	}

	// Tạo tên tệp tin độc bản bằng nano giây tránh trùng lặp
	ext := filepath.Ext(handler.Filename)
	uniqueFilename := fmt.Sprintf("img_%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, uniqueFilename)

	dst, err := os.Create(filePath)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Lỗi hệ thống khi tạo tệp tin ảnh")
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		h.writeJSONError(w, http.StatusInternalServerError, "Lỗi ghi dữ liệu ảnh")
		return
	}

	publicUrl := "/assets/uploads/" + uniqueFilename

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Tải lên ảnh thành công!",
		"url":     publicUrl,
	})
}

func (h *AdminHandler) writeJSONError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{
		"success": false,
		"error":   message,
	})
}

// ManageClass xử lý POST /api/v1/admin/classes/manage (Admin CRUD khóa học)
func (h *AdminHandler) ManageClass(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	var req struct {
		Action    string        `json:"action"` // "add", "update", "delete"
		ID        string        `json:"id"`
		ClassData *domain.Class `json:"class_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ")
		return
	}

	if req.Action == "add" {
		if req.ClassData == nil {
			h.writeJSONError(w, http.StatusBadRequest, "Thiếu thông tin lớp học cần thêm")
			return
		}
		// Generate dynamic standard ID if empty
		if req.ClassData.ID == "" {
			req.ClassData.ID = "class_" + fmt.Sprintf("%d", time.Now().UnixNano())
		}
		err := h.uc.AddClass(req.ClassData)
		if err != nil {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Thêm lớp học mới thành công!",
			"data":    req.ClassData,
		})
		return
	} else if req.Action == "update" {
		if req.ClassData == nil {
			h.writeJSONError(w, http.StatusBadRequest, "Thiếu thông tin lớp học cần sửa")
			return
		}
		err := h.uc.UpdateClass(req.ClassData)
		if err != nil {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Cập nhật lớp học thành công!",
		})
		return
	} else if req.Action == "delete" {
		if req.ID == "" {
			h.writeJSONError(w, http.StatusBadRequest, "Thiếu mã ID lớp học cần xóa")
			return
		}
		err := h.uc.DeleteClass(req.ID)
		if err != nil {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Xóa lớp học thành công khỏi hệ thống!",
		})
		return
	} else {
		h.writeJSONError(w, http.StatusBadRequest, "Thao tác không được hỗ trợ (chỉ chấp nhận 'add', 'update' hoặc 'delete')")
	}
}

// ManageGallery xử lý POST /api/v1/admin/gallery/manage
func (h *AdminHandler) ManageGallery(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeJSONError(w, http.StatusMethodNotAllowed, "Phương thức HTTP không được hỗ trợ")
		return
	}

	if session := h.authenticateAdmin(w, r); session == nil {
		return
	}

	var req struct {
		Action    string               `json:"action"` // "add" hoặc "delete"
		ID        int64                `json:"id"`
		ImageData *domain.GalleryImage `json:"image_data"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.writeJSONError(w, http.StatusBadRequest, "Dữ liệu yêu cầu không hợp lệ")
		return
	}

	if req.Action == "add" {
		if req.ImageData == nil {
			h.writeJSONError(w, http.StatusBadRequest, "Thiếu dữ liệu hình ảnh cần thêm")
			return
		}

		err := h.uc.AddGalleryImage(req.ImageData)
		if err != nil {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Thêm hình ảnh hoạt động lớp học thành công!",
			"data":    req.ImageData,
		})
		return
	} else if req.Action == "delete" {
		err := h.uc.DeleteGalleryImage(req.ID)
		if err != nil {
			h.writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Xóa hình ảnh hoạt động lớp học thành công!",
		})
		return
	} else {
		h.writeJSONError(w, http.StatusBadRequest, "Hành động (Action) không hợp lệ")
	}
}
