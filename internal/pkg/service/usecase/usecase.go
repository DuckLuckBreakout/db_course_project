package usecase

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/service"
)

type UseCase struct {
	Repository service.Repository
}

func (u UseCase) Close() {
	u.Repository.Close()
}

func (u UseCase) Clear() error {
	return u.Repository.Clear()
}

func (u UseCase) Status() (*models.Status, error) {
	return u.Repository.Status()
}

func NewUseCase(repo service.Repository) service.UseCase {
	return &UseCase{
		Repository: repo,
	}
}
