package usecase

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread"
)

type UseCase struct {
	Repository thread.Repository
}

func (u UseCase) Details(thread *models.Thread) error {
	return u.Repository.Details(thread)
}

func (u UseCase) Vote(thread *models.ThreadVoice) (*models.Thread, error) {
	return u.Repository.Vote(thread)
}

func NewUseCase(repo thread.Repository) thread.UseCase {
	return &UseCase{
		Repository: repo,
	}
}
