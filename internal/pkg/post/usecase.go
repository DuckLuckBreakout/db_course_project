package post

import "github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"

type UseCase interface {
	Details(id int) (*models.Post, error)
	DetailsUser(id int) (*models.User, error)
	DetailsForum(id int) (*models.Forum, error)
	DetailsThread(id int) (*models.Thread, error)
	UpdateDetails(updatePost *models.Post) (*models.Post, error)
	Close()
}
