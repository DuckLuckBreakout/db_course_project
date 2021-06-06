package thread

import "github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"

type Repository interface {
	Vote(thread *models.ThreadVoice) (*models.Thread, error)
}
