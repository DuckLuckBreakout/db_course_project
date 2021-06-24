package user

import "github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"

type UseCase interface {
	Create(user *models.User) ([]*models.User, error)
	Profile(user *models.User) error
	UpdateProfile(user *models.User) error
}
