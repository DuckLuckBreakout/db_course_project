package usecase

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum"
)

type UseCase struct {
	Repository forum.Repository
}

func NewUseCase(repo forum.Repository) forum.UseCase {
	return &UseCase{
		Repository: repo,
	}
}
