package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Khai báo các hằng số nghiệp vụ
const (
	ModelOnline  = "online"
	ModelOffline = "offline"

	TypeBasic       = "basic"
	TypeAdvanced    = "advanced"
	TypeHighQuality = "high_quality"
)

// Lead đại diện cho thực thể Phụ huynh đăng ký học thử của hệ thống
type Lead struct {
	ID                  int64     `json:"id"`
	ParentName          string    `json:"parent_name"`
	PhoneNumber         string    `json:"phone_number"`
	StudentName         string    `json:"student_name"`
	Grade               int       `json:"grade"`
	LearningModel       string    `json:"learning_model"`
	ClassType           string    `json:"class_type"`
	AcademicPerformance string    `json:"academic_performance"` // "excellent" (Giỏi/Xuất sắc), "good" (Khá), "average" (Trung bình)
	ConsultedBy         string    `json:"consulted_by"`         // Lưu tên tài khoản đã tư vấn
	Status              string    `json:"status"`               // "pending", "contacted", "enrolled", "cancelled"
	CreatedAt           time.Time `json:"created_at"`
}

// LeadRepository định nghĩa interface tương tác dữ liệu cho Lead
type LeadRepository interface {
	Save(lead *Lead) error
	FindAll() ([]*Lead, error)
	UpdateStatus(id int64, status string, consultedBy string) error
}


// Validate là hàm kiểm tra luật nghiệp vụ cốt lõi tại tầng Domain
func (l *Lead) Validate() error {
	// 1. Kiểm tra các trường bắt buộc không được trống
	l.ParentName = strings.TrimSpace(l.ParentName)
	l.StudentName = strings.TrimSpace(l.StudentName)
	l.PhoneNumber = strings.TrimSpace(l.PhoneNumber)

	if l.ParentName == "" {
		return errors.New("họ và tên Phụ huynh không được để trống")
	}
	if l.StudentName == "" {
		return errors.New("họ và tên Học sinh không được để trống")
	}
	if l.PhoneNumber == "" {
		return errors.New("số điện thoại liên hệ không được để trống")
	}

	// 2. Kiểm tra định dạng số điện thoại Việt Nam (10 số, bắt đầu bằng 03, 05, 07, 08, 09)
	phoneRegex := regexp.MustCompile(`^(03|05|07|08|09)\d{8}$`)
	if !phoneRegex.MatchString(l.PhoneNumber) {
		return fmt.Errorf("số điện thoại '%s' không đúng định dạng Việt Nam (phải gồm 10 chữ số và bắt đầu bằng các đầu số di động hợp lệ)", l.PhoneNumber)
	}

	// 3. Kiểm tra khối lớp hợp lệ (Lớp 1 đến Lớp 9)
	if l.Grade < 1 || l.Grade > 9 {
		return fmt.Errorf("khối lớp phải nằm trong khoảng từ lớp 1 đến lớp 9 (nhận được: Lớp %d)", l.Grade)
	}

	// 4. Kiểm tra hình thức học hợp lệ
	l.LearningModel = strings.ToLower(strings.TrimSpace(l.LearningModel))
	if l.LearningModel != ModelOnline && l.LearningModel != ModelOffline {
		return fmt.Errorf("hình thức học không hợp lệ, phải là 'online' hoặc 'offline' (nhận được: '%s')", l.LearningModel)
	}

	// 5. Kiểm tra loại lớp học hợp lệ
	l.ClassType = strings.ToLower(strings.TrimSpace(l.ClassType))
	if l.ClassType != TypeBasic && l.ClassType != TypeAdvanced && l.ClassType != TypeHighQuality {
		return fmt.Errorf("loại lớp không hợp lệ, phải là 'basic', 'advanced' hoặc 'high_quality' (nhận được: '%s')", l.ClassType)
	}

	// 6. 🔥 LUẬT NGHIỆP VỤ ĐẶC BIỆT: Lớp chất lượng cao chỉ dành riêng cho Lớp 5 ôn thi chuyên cấp 2
	if l.ClassType == TypeHighQuality && l.Grade != 5 {
		return fmt.Errorf("không thể đăng ký lớp Chất lượng cao (CLC / Luyện thi chuyên) cho học sinh Lớp %d. Lớp CLC chỉ được mở duy nhất cho khối Lớp 5 ôn thi vào cấp 2", l.Grade)
	}

	// 7. Kiểm tra học lực hiện tại hợp lệ
	l.AcademicPerformance = strings.ToLower(strings.TrimSpace(l.AcademicPerformance))
	if l.AcademicPerformance == "" {
		l.AcademicPerformance = "good" // Mặc định là Khá nếu trống
	}
	if l.AcademicPerformance != "excellent" && l.AcademicPerformance != "good" && l.AcademicPerformance != "average" {
		return fmt.Errorf("học lực học tập không hợp lệ (nhận được: '%s')", l.AcademicPerformance)
	}

	return nil
}
