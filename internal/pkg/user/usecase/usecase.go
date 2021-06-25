package usecase

import (
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/user"
)

type UseCase struct {
	Repository user.Repository
}

func NewUseCase(repo user.Repository) user.UseCase {
	return &UseCase{
		Repository: repo,
	}
}

func (u UseCase) Create(user *models.User) ([]*models.User, error) {
	if err := u.Repository.Create(user); err != nil {
		users, err := u.Repository.GetAllUsersByNicknameAndEmail(user)
		if err != nil {
			return nil, err
		}
		return users, errors.ErrUserAlreadyCreatedError
	}
	result := make([]*models.User, 0)
	result = append(result, user)
	return result, nil
}

func (u UseCase) Profile(user *models.User) error {
	if err := u.Repository.GetUserByNickname(user); err != nil {
		return err
	}
	return nil
}

func (u UseCase) UpdateProfile(user *models.User) error {
	if err := u.Repository.Update(user); err != nil {
		return err
	}
	return nil
}
