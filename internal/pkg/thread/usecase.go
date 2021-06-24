package thread

import "github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"

type UseCase interface {
	Vote(thread *models.ThreadVoice) (*models.Thread, error)
	Details(thread *models.Thread) error
	Posts(thread *models.PostSearch) ([]*models.Post, error)
	UpdateDetails(thread *models.ThreadUpdate) (*models.Thread, error)
	Create(slugOrId string, posts []*models.Post) error
	Close()
}
