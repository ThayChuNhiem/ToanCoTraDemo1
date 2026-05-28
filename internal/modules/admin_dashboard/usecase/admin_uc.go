package usecase

import (
	"errors"
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/internal/repository"
)

type AdminUsecase struct {
	repo *repository.CentralRepository
}

func NewAdminUsecase(repo *repository.CentralRepository) *AdminUsecase {
	return &AdminUsecase{repo: repo}
}

// GetSalesReport tổng hợp số liệu thống kê trực quan
func (uc *AdminUsecase) GetSalesReport() (*domain.SalesReport, error) {
	leads, err := uc.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var total, pending, contacted int64
	var online, offline int
	gradeDist := make(map[int]int)

	for _, l := range leads {
		total++
		if l.Status == "contacted" {
			contacted++
		} else {
			pending++
		}

		if l.LearningModel == "online" {
			online++
		} else {
			offline++
		}

		gradeDist[l.Grade]++
	}

	var onlinePct, offlinePct float64
	if total > 0 {
		onlinePct = (float64(online) / float64(total)) * 100
		offlinePct = (float64(offline) / float64(total)) * 100
	}

	return &domain.SalesReport{
		TotalLeads:        total,
		ContactedLeads:    contacted,
		PendingLeads:      pending,
		OnlinePercentage:  onlinePct,
		OfflinePercentage: offlinePct,
		GradeDistribution: gradeDist,
	}, nil
}

// UpdateClassPrice cập nhật học phí lớp học
func (uc *AdminUsecase) UpdateClassPrice(id string, price string) error {
	if id == "" || price == "" {
		return errors.New("mã lớp học và học phí mới không được để trống")
	}
	return uc.repo.UpdatePrice(id, price)
}

// UpdateTeacherBio cập nhật tiểu sử giáo viên
func (uc *AdminUsecase) UpdateTeacherBio(teacher *domain.Teacher) error {
	if teacher.ID <= 0 || teacher.Name == "" || teacher.Bio == "" {
		return errors.New("thông tin ID, tên giáo viên và tiểu sử không được để trống")
	}
	return uc.repo.UpdateTeacher(teacher)
}

// AddHonoredStudent thêm học sinh bảng vàng mới
func (uc *AdminUsecase) AddHonoredStudent(student *domain.Student) error {
	if student.Name == "" || student.Class == "" || student.Achievement == "" {
		return errors.New("họ tên học sinh, lớp học và thành tích xuất sắc không được để trống")
	}
	return uc.repo.SaveStudent(student)
}

// DeleteHonoredStudent xóa học sinh bảng vàng
func (uc *AdminUsecase) DeleteHonoredStudent(id int64) error {
	if id <= 0 {
		return errors.New("mã ID học sinh vinh danh không hợp lệ")
	}
	return uc.repo.DeleteStudent(id)
}

// CreateTeacherAccount tạo tài khoản giáo viên/sales mới (CHỈ ADMIN CÓ QUYỀN)
func (uc *AdminUsecase) CreateTeacherAccount(user *domain.User) error {
	if user.Username == "" || user.Password == "" || user.Name == "" {
		return errors.New("tên tài khoản, mật khẩu và họ tên giáo viên không được để trống")
	}
	
	// Ép buộc vai trò là teacher/sale
	user.Role = domain.RoleTeacher
	return uc.repo.SaveUser(user)
}

// AddTeacher thêm giáo viên/cộng sự mới
func (uc *AdminUsecase) AddTeacher(teacher *domain.Teacher) error {
	if teacher.Name == "" || teacher.Role == "" || teacher.Bio == "" {
		return errors.New("tên giáo viên, vai trò giảng dạy và tiểu sử không được để trống")
	}
	return uc.repo.SaveTeacher(teacher)
}

// DeleteTeacher xóa giáo viên/cộng sự khỏi hệ thống
func (uc *AdminUsecase) DeleteTeacher(id int64) error {
	if id <= 0 {
		return errors.New("mã ID giáo viên không hợp lệ")
	}
	return uc.repo.DeleteTeacher(id)
}

// GetAllUsers lấy toàn bộ danh sách tài khoản
func (uc *AdminUsecase) GetAllUsers() ([]*domain.User, error) {
	return uc.repo.FindAllUsers()
}

// DeleteSalesUser xóa tài khoản sales
func (uc *AdminUsecase) DeleteSalesUser(username string) error {
	if username == "" {
		return errors.New("tên tài khoản cần xóa không được để trống")
	}
	return uc.repo.DeleteUser(username)
}

// AddClass thêm khóa học mới
func (uc *AdminUsecase) AddClass(class *domain.Class) error {
	if class.ID == "" || class.Name == "" || class.Price == "" {
		return errors.New("mã lớp học, tên lớp học và học phí không được để trống")
	}
	return uc.repo.SaveClass(class)
}

// UpdateClass cập nhật thông tin khóa học toàn diện
func (uc *AdminUsecase) UpdateClass(class *domain.Class) error {
	if class.ID == "" || class.Name == "" || class.Price == "" {
		return errors.New("mã lớp học, tên lớp học và học phí không được để trống")
	}
	return uc.repo.UpdateClass(class)
}

// DeleteClass xóa lớp học khỏi hệ thống
func (uc *AdminUsecase) DeleteClass(id string) error {
	if id == "" {
		return errors.New("mã lớp học không hợp lệ")
	}
	return uc.repo.DeleteClass(id)
}

