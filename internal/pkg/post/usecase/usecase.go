package usecase

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/post"
)

type UseCase struct {
	Repository post.Repository
}

func (u UseCase) DetailsUser(id int) (*models.User, error) {
	return u.Repository.DetailsUser(id)
}

func (u UseCase) DetailsForum(id int) (*models.Forum, error) {
	return u.Repository.DetailsForum(id)
}

func (u UseCase) DetailsThread(id int) (*models.Thread, error) {
	return u.Repository.DetailsThread(id)
}

func (u UseCase) UpdateDetails(updatePost *models.Post) (*models.Post, error) {
	return u.Repository.UpdateDetails(updatePost)
}

func (u UseCase) Details(id int) (*models.Post, error) {
	return u.Repository.Details(id)
}

func NewUseCase(repo post.Repository) post.UseCase {
	return &UseCase{
		Repository: repo,
	}
}
