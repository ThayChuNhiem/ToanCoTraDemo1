package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sort"
	"sync"
	"time"
	"toan-co-tra-backend/internal/domain"
)

// DataStore là cấu trúc dữ liệu lưu trữ toàn bộ trong database_mock.json
type DataStore struct {
	Leads    []*domain.Lead          `json:"leads"`
	Classes  []*domain.Class         `json:"classes"`
	Teachers []*domain.Teacher       `json:"teachers"`
	Students []*domain.Student       `json:"students"`
	Users    []*domain.User          `json:"users"`
	Gallery  []*domain.GalleryImage  `json:"gallery"` // [MỚI] Hình ảnh hoạt động lớp học
}

// CentralRepository là trình quản lý kho dữ liệu tập trung an toàn đa luồng,
// tự động tải và lưu trữ trạng thái xuống tệp tin database_mock.json.
type CentralRepository struct {
	mu       sync.RWMutex
	filePath string
	store    *DataStore
	nextLeadID    int64
	nextTeacherID int64
	nextStudentID int64
	nextGalleryID int64 // [MỚI] ID tự tăng của hình ảnh gallery
}

// NewCentralRepository khởi tạo CentralRepository
func NewCentralRepository(filePath string) *CentralRepository {
	repo := &CentralRepository{
		filePath: filePath,
		store: &DataStore{
			Leads:    make([]*domain.Lead, 0),
			Classes:  make([]*domain.Class, 0),
			Teachers: make([]*domain.Teacher, 0),
			Students: make([]*domain.Student, 0),
			Users:    make([]*domain.User, 0),
			Gallery:  make([]*domain.GalleryImage, 0),
		},
		nextLeadID:    1,
		nextTeacherID: 1,
		nextStudentID: 1,
		nextGalleryID: 1,
	}

	// Tải dữ liệu hoặc khởi tạo dữ liệu mẫu nếu chưa có tệp tin
	if err := repo.load(); err != nil {
		fmt.Printf("[WARN] Không thể tải database_mock.json (%s), tiến hành gieo hạt dữ liệu mẫu mặc định...\n", err.Error())
		repo.seedDefaultData()
		_ = repo.save()
	} else {
		// Nếu tải thành công nhưng trường gallery trống (ví dụ do database_mock.json cũ chưa có gallery),
		// gieo hạt một số hình ảnh lớp học thực tế đẹp mắt để UI hiển thị lộng lẫy ngay lập tức.
		if len(repo.store.Gallery) == 0 {
			repo.store.Gallery = []*domain.GalleryImage{
				{
					ID:  1,
					URL: "https://images.unsplash.com/photo-1577896851231-70ef18881754?q=80&w=600&auto=format&fit=crop",
				},
				{
					ID:  2,
					URL: "https://images.unsplash.com/photo-1427504494785-3a9ca7044f45?q=80&w=600&auto=format&fit=crop",
				},
				{
					ID:  3,
					URL: "https://images.unsplash.com/photo-1509062522246-3755977927d7?q=80&w=600&auto=format&fit=crop",
				},
			}
			repo.nextGalleryID = 4
			_ = repo.save()
		}
	}

	return repo
}

// load tải dữ liệu từ tệp tin JSON
func (r *CentralRepository) load() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	file, err := os.Open(r.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var ds DataStore
	if err := json.Unmarshal(bytes, &ds); err != nil {
		return err
	}

	r.store = &ds

	// Cập nhật các ID tự tăng lớn nhất
	for _, l := range r.store.Leads {
		if l.ID >= r.nextLeadID {
			r.nextLeadID = l.ID + 1
		}
	}
	for _, t := range r.store.Teachers {
		if t.ID >= r.nextTeacherID {
			r.nextTeacherID = t.ID + 1
		}
	}
	for _, s := range r.store.Students {
		if s.ID >= r.nextStudentID {
			r.nextStudentID = s.ID + 1
		}
	}
	for _, g := range r.store.Gallery {
		if g.ID >= r.nextGalleryID {
			r.nextGalleryID = g.ID + 1
		}
	}

	return nil
}

// save ghi lưu trữ dữ liệu hiện tại xuống tệp tin JSON
func (r *CentralRepository) save() error {
	// Lưu ý: Hàm này được gọi khi đã có Lock bên ngoài gọi nó
	bytes, err := json.MarshalIndent(r.store, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(r.filePath, bytes, 0644)
}

// seedDefaultData tạo dữ liệu hạt giống ban đầu lộng lẫy và hoàn chỉnh
func (r *CentralRepository) seedDefaultData() {
	// Gieo hạt hệ thống lớp học động từ lớp 1 đến lớp 9
	r.store.Classes = make([]*domain.Class, 0)
	for grade := 1; grade <= 9; grade++ {
		// 1. Cơ Bản Online
		r.store.Classes = append(r.store.Classes, &domain.Class{
			ID:        fmt.Sprintf("class_%d_basic_online", grade),
			Name:      fmt.Sprintf("Toán Tư Duy Cơ Bản Lớp %d (Online)", grade),
			Grade:     grade,
			Type:      "basic",
			Model:     "online",
			Price:     "800.000đ/tháng",
			Desc:      fmt.Sprintf("Lớp học Toán tư duy cơ bản trực tuyến Lớp %d, bám sát chương trình của Bộ Giáo Dục, tạo nền tảng vững vàng.", grade),
			IsPopular: false,
		})

		// 2. Cơ Bản Offline
		r.store.Classes = append(r.store.Classes, &domain.Class{
			ID:        fmt.Sprintf("class_%d_basic_offline", grade),
			Name:      fmt.Sprintf("Toán Tư Duy Cơ Bản Lớp %d (Tại Lớp)", grade),
			Grade:     grade,
			Type:      "basic",
			Model:     "offline",
			Price:     "1.000.000đ/tháng",
			Desc:      fmt.Sprintf("Lớp học trực tiếp tại trung tâm, rèn luyện kỹ năng giải toán cơ bản và tư duy logic Lớp %d.", grade),
			IsPopular: false,
		})

		// 3. Nâng Cao Online
		r.store.Classes = append(r.store.Classes, &domain.Class{
			ID:        fmt.Sprintf("class_%d_adv_online", grade),
			Name:      fmt.Sprintf("Toán Tư Duy Nâng Cao Lớp %d (Online)", grade),
			Grade:     grade,
			Type:      "advanced",
			Model:     "online",
			Price:     "1.000.000đ/tháng",
			Desc:      fmt.Sprintf("Học trực tuyến nâng cao Lớp %d, rèn luyện tư duy sâu sắc, chinh phục bài toán logic phức tạp.", grade),
			IsPopular: false,
		})

		// 4. Nâng Cao Offline
		r.store.Classes = append(r.store.Classes, &domain.Class{
			ID:        fmt.Sprintf("class_%d_adv_offline", grade),
			Name:      fmt.Sprintf("Toán Tư Duy Nâng Cao Lớp %d (Tại Lớp)", grade),
			Grade:     grade,
			Type:      "advanced",
			Model:     "offline",
			Price:     "1.200.000đ/tháng",
			Desc:      fmt.Sprintf("Lớp học bồi dưỡng nâng cao Lớp %d trực tiếp tại trung tâm, chinh phục các kì thi HSG.", grade),
			IsPopular: grade == 5, // Lớp 5 Nâng Cao mặc định hot
		})

		// 5. Riêng Lớp 5 có hệ Chất Lượng Cao ôn thi chuyên cấp 2
		if grade == 5 {
			r.store.Classes = append(r.store.Classes, &domain.Class{
				ID:        "class_5_clc_offline",
				Name:      "Toán Chuyên Lớp 5 CLC (Tại Lớp)",
				Grade:     5,
				Type:      "high_quality",
				Model:     "offline",
				Price:     "1.800.000đ/tháng",
				Desc:      "Lộ trình tinh gọn chuyên biệt ôn thi bứt phá trực tiếp tại trung tâm vào các trường THCS chuyên & CLC tại Hà Nội.",
				IsPopular: true,
			})
			r.store.Classes = append(r.store.Classes, &domain.Class{
				ID:        "class_5_clc_online",
				Name:      "Toán Chuyên Lớp 5 CLC (Online)",
				Grade:     5,
				Type:      "high_quality",
				Model:     "online",
				Price:     "1.500.000đ/tháng",
				Desc:      "Lộ trình ôn thi chuyên trực tuyến tương tác cao, bứt phá thi chuyên dành riêng cho học sinh khối lớp 5.",
				IsPopular: false,
			})
		}
	}

	// Gieo hạt giảng viên (Cô Trà và cộng sự)
	r.store.Teachers = []*domain.Teacher{
		{
			ID:        1,
			Name:      "Cô Giáo Trà (TCT)",
			Role:      "Giáo viên Sáng lập hệ thống",
			Avatar:    "/assets/cotra.jpg",
			Bio:       "Tốt nghiệp xuất sắc chuyên ngành Sư phạm Toán học. Với hơn 8 năm kinh nghiệm giảng dạy Toán tiểu học và luyện thi chuyên cấp 2. Cô Trà đã dẫn dắt hơn 500 học sinh đỗ vào các trường điểm Hà Nội - Amsterdam, Thanh Xuân, Cầu Giấy. Châm ngôn: Ươm mầm tình yêu Toán học để khơi dậy tiềm năng tự thân của mỗi đứa trẻ.",
			Education: "Thạc sĩ Khoa học Giáo dục - Đại học Sư Phạm Hà Nội",
			Order:     1,
		},
		{
			ID:        2,
			Name:      "Thầy Nguyễn Đức Anh",
			Role:      "Cố vấn chuyên môn sư phạm",
			Avatar:    "",
			Bio:       "Cựu học sinh chuyên Toán ĐHQG Hà Nội. Thầy có chuyên môn cao trong việc biên soạn hệ thống giáo trình Toán tư duy liên cấp bám sát đề thi chuyên cấp 2 và cấp 3. Thầy phụ trách giảng dạy các lớp CLC mũi nhọn.",
			Education: "Cử nhân Toán Tin - Đại học Bách Khoa Hà Nội",
			Order:     2,
		},
		{
			ID:        3,
			Name:      "Cô Phạm Minh Châu",
			Role:      "Giáo viên Toán Tư duy tiểu học",
			Avatar:    "",
			Bio:       "Nhiệt huyết, kiên nhẫn và vô cùng sát sao hỗ trợ các con. Cô có phương pháp trực quan hóa sinh động giúp các con mất gốc tự tin bứt phá, xóa bỏ nỗi sợ Toán và rèn kỹ năng tính toán siêu tốc.",
			Education: "Cử nhân Giáo dục Tiểu học - Đại học Sư Phạm Hà Nội",
			Order:     3,
		},
	}

	// Gieo hạt học sinh vinh danh
	r.store.Students = []*domain.Student{
		{
			ID:          1,
			Name:        "Nguyễn Minh Đức",
			Class:       "Lớp 5 CLC Chuyên",
			Year:        "Niên khóa 2024 - 2025",
			Achievement: "Đỗ Chuyên Toán THCS chuyên Hà Nội - Amsterdam (Thủ khoa môn Toán)",
			Avatar:      "",
			Order:       1,
		},
		{
			ID:          2,
			Name:        "Phạm Khánh An",
			Class:       "Lớp 5 Nâng Cao",
			Year:        "Niên khóa 2024 - 2025",
			Achievement: "Huy chương Vàng Olympic Toán Quốc tế TIMO & SASMO",
			Avatar:      "",
			Order:       2,
		},
		{
			ID:          3,
			Name:        "Lê Bảo Nam",
			Class:       "Lớp 5 CLC Chuyên",
			Year:        "Niên khóa 2024 - 2025",
			Achievement: "Đỗ lớp chọn CLC THCS Cầu Giấy & THCS Thanh Xuân (Điểm Toán 9.75)",
			Avatar:      "",
			Order:       3,
		},
	}

	// Gieo hạt tài khoản hệ thống (admin1 & teacher1)
	r.store.Users = []*domain.User{
		{
			Username: "admin1",
			Password: "admin123",
			Name:     "Quản trị viên tối cao",
			Role:     domain.RoleAdmin,
		},
		{
			Username: "teacher1",
			Password: "password123",
			Name:     "Cô Giáo Trà",
			Role:     domain.RoleTeacher,
		},
	}

	// Gieo hạt Leads ban đầu cho Sales feed sinh động
	r.store.Leads = []*domain.Lead{
		{
			ID:                  1,
			ParentName:          "Nguyễn Thị Lan",
			PhoneNumber:         "0983459912",
			StudentName:         "Nguyễn Minh Đức",
			Grade:               5,
			LearningModel:       "online",
			ClassType:           "basic",
			AcademicPerformance: "good",
			Status:              "pending",
			CreatedAt:           time.Now().Add(-10 * time.Minute),
		},
		{
			ID:                  2,
			ParentName:          "Trần Văn Hùng",
			PhoneNumber:         "0976543210",
			StudentName:         "Trần Khánh An",
			Grade:               4,
			LearningModel:       "offline",
			ClassType:           "advanced",
			AcademicPerformance: "excellent",
			Status:              "contacted",
			CreatedAt:           time.Now().Add(-30 * time.Minute),
		},
	}

	// Gieo hạt hình ảnh lớp học
	r.store.Gallery = []*domain.GalleryImage{
		{
			ID:  1,
			URL: "/assets/banner.png",
		},
	}

	r.nextLeadID = 3
	r.nextTeacherID = 4
	r.nextStudentID = 4
	r.nextGalleryID = 2
}

// -----------------------------------------------------------------------------
// IMPLEMENTATION OF domain.LeadRepository
// -----------------------------------------------------------------------------

func (r *CentralRepository) Save(lead *domain.Lead) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	lead.ID = r.nextLeadID
	r.nextLeadID++
	lead.Status = "pending"
	lead.CreatedAt = time.Now()

	r.store.Leads = append(r.store.Leads, lead)
	return r.save()
}

func (r *CentralRepository) FindAll() ([]*domain.Lead, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copied := make([]*domain.Lead, len(r.store.Leads))
	for i, l := range r.store.Leads {
		// Bản sao sâu nông vừa đủ
		copied[i] = &domain.Lead{
			ID:                  l.ID,
			ParentName:          l.ParentName,
			PhoneNumber:         l.PhoneNumber,
			StudentName:         l.StudentName,
			Grade:               l.Grade,
			LearningModel:       l.LearningModel,
			ClassType:           l.ClassType,
			AcademicPerformance: l.AcademicPerformance,
			ConsultedBy:         l.ConsultedBy,
			Status:              l.Status,
			CreatedAt:           l.CreatedAt,
		}
	}
	return copied, nil
}

func (r *CentralRepository) UpdateStatus(id int64, status string, consultedBy string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, lead := range r.store.Leads {
		if lead.ID == id {
			lead.Status = status
			lead.ConsultedBy = consultedBy
			if status == "pending" {
				lead.ConsultedBy = ""
			}
			return r.save()
		}
	}
	return errors.New("không tìm thấy Lead học sinh với ID tương ứng")
}

// -----------------------------------------------------------------------------
// IMPLEMENTATION OF domain.ClassRepository
// -----------------------------------------------------------------------------

func (r *CentralRepository) FindAllClasses() ([]*domain.Class, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copied := make([]*domain.Class, len(r.store.Classes))
	for i, c := range r.store.Classes {
		copied[i] = &domain.Class{
			ID:        c.ID,
			Name:      c.Name,
			Grade:     c.Grade,
			Type:      c.Type,
			Model:     c.Model,
			Price:     c.Price,
			Desc:      c.Desc,
			IsPopular: c.IsPopular,
		}
	}
	return copied, nil
}

func (r *CentralRepository) UpdatePrice(id string, price string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, c := range r.store.Classes {
		if c.ID == id {
			c.Price = price
			return r.save()
		}
	}
	return errors.New("không tìm thấy lớp học yêu cầu")
}

// -----------------------------------------------------------------------------
// IMPLEMENTATION OF domain.TeacherRepository
// -----------------------------------------------------------------------------

func (r *CentralRepository) FindAllTeachers() ([]*domain.Teacher, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copied := make([]*domain.Teacher, len(r.store.Teachers))
	for i, t := range r.store.Teachers {
		copied[i] = &domain.Teacher{
			ID:        t.ID,
			Name:      t.Name,
			Role:      t.Role,
			Avatar:    t.Avatar,
			Bio:       t.Bio,
			Education: t.Education,
			Order:     t.Order,
		}
	}
	return copied, nil
}

func (r *CentralRepository) SaveTeacher(teacher *domain.Teacher) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	teacher.ID = r.nextTeacherID
	r.nextTeacherID++
	r.store.Teachers = append(r.store.Teachers, teacher)
	return r.save()
}

func (r *CentralRepository) UpdateTeacher(teacher *domain.Teacher) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, t := range r.store.Teachers {
		if t.ID == teacher.ID {
			t.Name = teacher.Name
			t.Role = teacher.Role
			if teacher.Avatar != "" {
				t.Avatar = teacher.Avatar
			}
			t.Bio = teacher.Bio
			t.Education = teacher.Education
			t.Order = teacher.Order
			return r.save()
		}
	}
	return errors.New("không tìm thấy thông tin giáo viên tương ứng")
}

func (r *CentralRepository) DeleteTeacher(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if id == 1 {
		return errors.New("không thể xóa tài khoản sáng lập của Cô Giáo Trà")
	}

	index := -1
	for i, t := range r.store.Teachers {
		if t.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("không tìm thấy thông tin giáo viên yêu cầu xóa")
	}

	r.store.Teachers = append(r.store.Teachers[:index], r.store.Teachers[index+1:]...)
	return r.save()
}


// -----------------------------------------------------------------------------
// IMPLEMENTATION OF domain.StudentRepository
// -----------------------------------------------------------------------------

func (r *CentralRepository) FindAllStudents() ([]*domain.Student, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copied := make([]*domain.Student, len(r.store.Students))
	for i, s := range r.store.Students {
		copied[i] = &domain.Student{
			ID:          s.ID,
			Name:        s.Name,
			Class:       s.Class,
			Year:        s.Year,
			Achievement: s.Achievement,
			Avatar:      s.Avatar,
			Order:       s.Order,
		}
	}

	// Sắp xếp danh sách học sinh theo thuộc tính Order tăng dần (số nhỏ xếp trước)
	sort.Slice(copied, func(i, j int) bool {
		if copied[i].Order != copied[j].Order {
			return copied[i].Order < copied[j].Order
		}
		return copied[i].ID < copied[j].ID
	})

	return copied, nil
}

func (r *CentralRepository) SaveStudent(student *domain.Student) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	student.ID = r.nextStudentID
	r.nextStudentID++
	r.store.Students = append(r.store.Students, student)
	return r.save()
}

func (r *CentralRepository) DeleteStudent(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	index := -1
	for i, s := range r.store.Students {
		if s.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("không tìm thấy thông tin học sinh vinh danh")
	}

	r.store.Students = append(r.store.Students[:index], r.store.Students[index+1:]...)
	return r.save()
}

// -----------------------------------------------------------------------------
// IMPLEMENTATION OF domain.UserRepository
// -----------------------------------------------------------------------------

func (r *CentralRepository) FindUser(username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, u := range r.store.Users {
		if u.Username == username {
			return &domain.User{
				Username: u.Username,
				Password: u.Password,
				Name:     u.Name,
				Role:     u.Role,
			}, nil
		}
	}
	return nil, errors.New("tên tài khoản đăng nhập không tồn tại")
}

func (r *CentralRepository) SaveUser(user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Kiểm tra xem user đã tồn tại chưa
	for _, u := range r.store.Users {
		if u.Username == user.Username {
			return errors.New("tên tài khoản đăng nhập đã tồn tại trên hệ thống")
		}
	}

	r.store.Users = append(r.store.Users, user)
	return r.save()
}

func (r *CentralRepository) FindAllUsers() ([]*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copied := make([]*domain.User, len(r.store.Users))
	for i, u := range r.store.Users {
		copied[i] = &domain.User{
			Username: u.Username,
			Name:     u.Name,
			Role:     u.Role,
		}
	}
	return copied, nil
}

func (r *CentralRepository) DeleteUser(username string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if username == "admin1" {
		return errors.New("không thể xóa tài khoản quản trị tối cao của Admin")
	}

	index := -1
	for i, u := range r.store.Users {
		if u.Username == username {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("không tìm thấy tài khoản yêu cầu xóa")
	}

	r.store.Users = append(r.store.Users[:index], r.store.Users[index+1:]...)
	return r.save()
}

func (r *CentralRepository) SaveClass(class *domain.Class) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, c := range r.store.Classes {
		if c.ID == class.ID {
			return errors.New("mã lớp học này đã tồn tại trên hệ thống")
		}
	}

	r.store.Classes = append(r.store.Classes, class)
	return r.save()
}

func (r *CentralRepository) UpdateClass(class *domain.Class) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, c := range r.store.Classes {
		if c.ID == class.ID {
			c.Name = class.Name
			c.Grade = class.Grade
			c.Type = class.Type
			c.Model = class.Model
			c.Price = class.Price
			c.Desc = class.Desc
			c.IsPopular = class.IsPopular
			return r.save()
		}
	}
	return errors.New("không tìm thấy thông tin lớp học để cập nhật")
}

func (r *CentralRepository) DeleteClass(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	index := -1
	for i, c := range r.store.Classes {
		if c.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("không tìm thấy lớp học yêu cầu xóa")
	}

	r.store.Classes = append(r.store.Classes[:index], r.store.Classes[index+1:]...)
	return r.save()
}

// -----------------------------------------------------------------------------
// IMPLEMENTATION OF domain.GalleryRepository
// -----------------------------------------------------------------------------

func (r *CentralRepository) FindAllImages() ([]*domain.GalleryImage, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	copied := make([]*domain.GalleryImage, len(r.store.Gallery))
	for i, g := range r.store.Gallery {
		copied[i] = &domain.GalleryImage{
			ID:  g.ID,
			URL: g.URL,
		}
	}
	return copied, nil
}

func (r *CentralRepository) SaveImage(img *domain.GalleryImage) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	img.ID = r.nextGalleryID
	r.nextGalleryID++
	r.store.Gallery = append(r.store.Gallery, img)
	return r.save()
}

func (r *CentralRepository) DeleteImage(id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	index := -1
	for i, g := range r.store.Gallery {
		if g.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New("không tìm thấy hình ảnh lớp học yêu cầu xóa")
	}

	r.store.Gallery = append(r.store.Gallery[:index], r.store.Gallery[index+1:]...)
	return r.save()
}

