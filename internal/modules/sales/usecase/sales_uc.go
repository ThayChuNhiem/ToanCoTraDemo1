package usecase

import (
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/internal/repository"
)

type SalesUsecase struct {
	repo *repository.CentralRepository
}

func NewSalesUsecase(repo *repository.CentralRepository) *SalesUsecase {
	return &SalesUsecase{repo: repo}
}

func (uc *SalesUsecase) GetAllLeads() ([]*domain.Lead, error) {
	return uc.repo.FindAll()
}

func (uc *SalesUsecase) UpdateLeadStatus(id int64, status string, consultedBy string) error {
	return uc.repo.UpdateStatus(id, status, consultedBy)
}
