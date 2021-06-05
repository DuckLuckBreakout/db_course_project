package forum

import "github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"

type UseCase interface {
	Create(forum *models.Forum) error
}
