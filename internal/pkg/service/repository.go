package service

import "github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"

type Repository interface {
	Clear() error
	Status() (*models.Status, error)
	Close()
}
