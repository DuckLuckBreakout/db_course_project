package usecase

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/user"
)

type UseCase struct {
	Repository   user.Repository
}

func NewUseCase(repo user.Repository) user.UseCase {
	return &UseCase{
		Repository:   repo,
	}
}

func (u UseCase) Create(user *models.User) ([]*models.User, error) {
	if err := u.Repository.Create(user); err != nil {
		users, err := u.Repository.GetAllUsersByNicknameAndEmail(user)
		if err != nil {
			return nil, nil
		}
		return users, nil
	}
	return nil, nil
}
