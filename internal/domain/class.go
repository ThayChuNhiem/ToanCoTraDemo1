package domain

// Class đại diện cho thực thể Lớp học của hệ thống
type Class struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Grade     int    `json:"grade"`
	Type      string `json:"type"`       // "basic", "advanced", "high_quality"
	Model     string `json:"model"`      // "online", "offline"
	Price     string `json:"price"`      // Giá học phí hiển thị
	Desc      string `json:"desc"`       // Mô tả lớp học
	IsPopular bool   `json:"is_popular"` // Gắn mác hot
}

// ClassRepository định nghĩa cách thức tương tác lưu trữ cho Lớp học
type ClassRepository interface {
	FindAllClasses() ([]*Class, error)
	UpdatePrice(id string, price string) error
	SaveClass(class *Class) error
	UpdateClass(class *Class) error
	DeleteClass(id string) error
}

