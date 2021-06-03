package user

import "github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"

type Repository interface {
	Create(user *models.User) error
	GetAllUsersByNicknameAndEmail(user *models.User) ([]*models.User, error)
	GetUserByNickname(user *models.User) error
}
