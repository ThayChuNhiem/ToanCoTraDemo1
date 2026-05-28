package usecase

import (
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/internal/repository"
)

type HonorUsecase struct {
	repo *repository.CentralRepository
}

func NewHonorUsecase(repo *repository.CentralRepository) *HonorUsecase {
	return &HonorUsecase{repo: repo}
}

func (uc *HonorUsecase) GetHonoredStudents() ([]*domain.Student, error) {
	return uc.repo.FindAllStudents()
}
