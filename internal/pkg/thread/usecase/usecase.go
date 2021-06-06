package usecase

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread"
)

type UseCase struct {
	Repository thread.Repository
}


func NewUseCase(repo thread.Repository) thread.UseCase {
	return &UseCase{
		Repository: repo,
	}
}
