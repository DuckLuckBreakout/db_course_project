package forum

import "github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"

type UseCase interface {
	Create(forum *models.Forum) error
	Details(forum *models.Forum) error
	CreateThread(thread *models.Thread) error
	Threads(thread *models.ThreadSearch) ([]*models.Thread, error)
	Users(searchParams *models.UserSearch) ([]*models.User, error)
}
