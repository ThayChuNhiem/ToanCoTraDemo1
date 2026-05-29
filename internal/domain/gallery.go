package domain

// GalleryImage đại diện cho hình ảnh lớp học thực tế tại trung tâm
type GalleryImage struct {
	ID  int64  `json:"id"`
	URL string `json:"url"` // Đường dẫn ảnh (avatar/assets hoặc upload)
}

// GalleryRepository định nghĩa cách lưu trữ hình ảnh hoạt động lớp học
type GalleryRepository interface {
	FindAllImages() ([]*GalleryImage, error)
	SaveImage(img *GalleryImage) error
	DeleteImage(id int64) error
}
