package domain

// Teacher đại diện cho thực thể giảng viên của hệ thống (Cô Trà & cộng sự)
type Teacher struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Role      string `json:"role"`      // "Giáo viên Sáng lập", "Giáo viên Toán tư duy", "Cố vấn sư phạm"
	Avatar    string `json:"avatar"`    // Ảnh đại diện dạng URL hoặc SVG/base64
	Bio       string `json:"bio"`       // Tiểu sử mô tả chi tiết quá trình học tập/giảng dạy
	Education string `json:"education"` // Trình độ học vấn (e.g., Cử nhân ĐH Sư Phạm Hà Nội)
	Order     int    `json:"order"`     // Thứ tự hiển thị trên website
}

// TeacherRepository định nghĩa cách thức tương tác lưu trữ cho Giáo viên
type TeacherRepository interface {
	FindAllTeachers() ([]*Teacher, error)
	SaveTeacher(teacher *Teacher) error
	UpdateTeacher(teacher *Teacher) error
	DeleteTeacher(id int64) error
}


