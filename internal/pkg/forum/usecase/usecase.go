package usecase

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
)

type UseCase struct {
	Repository forum.Repository
}

func (u UseCase) Users(searchParams *models.UserSearch) ([]*models.User, error) {
	return u.Repository.Users(searchParams)
}

func (u UseCase) Threads(thread *models.ThreadSearch, sinceString string) ([]*models.Thread, error) {
	return u.Repository.Threads(thread, sinceString)
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
