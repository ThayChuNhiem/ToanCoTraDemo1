package usecase

import (
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/internal/repository"
)

type HomepageUsecase struct {
	repo *repository.CentralRepository
}

func NewHomepageUsecase(repo *repository.CentralRepository) *HomepageUsecase {
	return &HomepageUsecase{repo: repo}
}

func (uc *HomepageUsecase) GetClasses() ([]*domain.Class, error) {
	return uc.repo.FindAllClasses()
}

func (uc *HomepageUsecase) RegisterLead(lead *domain.Lead) (*domain.Lead, error) {
	// 1. Thực hiện validate nghiệp vụ tại tầng Domain
	if err := lead.Validate(); err != nil {
		return nil, err
	}

	// 2. Lưu trữ dữ liệu tuyển sinh vào kho lưu trữ
	if err := uc.repo.Save(lead); err != nil {
		return nil, err
	}

	return lead, nil
}
