package usecase

import (
	"toan-co-tra-backend/internal/domain"
	"toan-co-tra-backend/internal/repository"
)

type AboutUsecase struct {
	repo *repository.CentralRepository
}

func NewAboutUsecase(repo *repository.CentralRepository) *AboutUsecase {
	return &AboutUsecase{repo: repo}
}

func (uc *AboutUsecase) GetTeachers() ([]*domain.Teacher, error) {
	return uc.repo.FindAllTeachers()
}
