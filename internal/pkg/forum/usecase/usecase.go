package usecase

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
)

type UseCase struct {
	Repository forum.Repository
}

func (u UseCase) CreateThread(thread *models.Thread) error {
	return u.Repository.CreateThread(thread)
}

func (u UseCase) Details(forum *models.Forum) error {
	return u.Repository.Details(forum)
}

func (u UseCase) Create(forum *models.Forum) error {
	return u.Repository.Create(forum)
}

func NewUseCase(repo forum.Repository) forum.UseCase {
	return &UseCase{
		Repository: repo,
	}
}
