package domain

// Student đại diện cho học sinh đạt thành tích xuất sắc được tuyên dương
type Student struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Class       string `json:"class"`       // Lớp học (e.g., Lớp 5 CLC, Lớp 9 chuyên Toán)
	Year        string `json:"year"`        // Niên khóa vinh danh (e.g., 2024 - 2025)
	Achievement string `json:"achievement"` // Thành tích đạt được (e.g., Đỗ chuyên Hà Nội - Amsterdam)
	Avatar      string `json:"avatar"`      // Ảnh chân dung học sinh
}

// StudentRepository định nghĩa cách thức tương tác lưu trữ cho Học sinh tuyên dương
type StudentRepository interface {
	FindAllStudents() ([]*Student, error)
	SaveStudent(student *Student) error
	DeleteStudent(id int64) error
}

