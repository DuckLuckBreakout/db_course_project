package usecase

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/service"
)

type UseCase struct {
	Repository service.Repository
}

func (u UseCase) Clear() error {
	return u.Repository.Clear()
}

func (u UseCase) Status() (*models.Status, error) {
	panic("implement me")
}

func NewUseCase(repo service.Repository) service.UseCase {
	return &UseCase{
		Repository: repo,
	}
}
