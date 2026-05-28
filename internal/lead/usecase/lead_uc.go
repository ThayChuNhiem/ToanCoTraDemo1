package usecase

import (
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/pkg/event"
)

// LeadUsecase đại diện cho nghiệp vụ quản lý lead đăng ký học thử
type LeadUsecase struct {
	repo       domain.LeadRepository
	dispatcher *event.EventDispatcher
}

// NewLeadUsecase khởi tạo một LeadUsecase mới
func NewLeadUsecase(repo domain.LeadRepository, dispatcher *event.EventDispatcher) *LeadUsecase {
	return &LeadUsecase{
		repo:       repo,
		dispatcher: dispatcher,
	}
}

// RegisterNewLead xử lý quy trình Phụ huynh gửi thông tin đăng ký học thử
func (uc *LeadUsecase) RegisterNewLead(lead *domain.Lead) (*domain.Lead, error) {
	// 1. Kích hoạt kiểm tra luật nghiệp vụ (Validation) tại tầng Domain
	if err := lead.Validate(); err != nil {
		return nil, err
	}

	// 2. Ghi nhận thông tin vào Cơ sở dữ liệu
	if err := uc.repo.Save(lead); err != nil {
		return nil, err
	}

	// 3. Phát sự kiện đăng ký thành công cho Live Monitor SSE và các module thông báo sau này
	uc.dispatcher.Publish("registration.created", lead)

	return lead, nil
}

// GetAllLeads lấy toàn bộ danh sách đăng ký tư vấn học sinh
func (uc *LeadUsecase) GetAllLeads() ([]*domain.Lead, error) {
	return uc.repo.FindAll()
}

// UpdateLeadStatus cập nhật trạng thái tư vấn (Đã tư vấn/Chưa tư vấn) của Lead
func (uc *LeadUsecase) UpdateLeadStatus(id int64, status string) error {
	return uc.repo.UpdateStatus(id, status, "")
}
