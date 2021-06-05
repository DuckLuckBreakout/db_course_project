package forum

import "github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"

type Repository interface {
	Create(forum *models.Forum) error
}
