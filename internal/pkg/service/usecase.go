package service

import "github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"

type UseCase interface {
	Clear() error
	Status() (*models.Status, error)
	Close()
}
