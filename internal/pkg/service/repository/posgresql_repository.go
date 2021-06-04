package repository

import (
	"database/sql"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/service"
)

type Repository struct {
	db *sql.DB
}

func (r Repository) Clear() error {
	row := r.db.QueryRow("TRUNCATE TABLE users")
	if err := row.Err(); err != nil {
		return err
	}
	return nil
}

func (r Repository) Status() (*models.Status, error) {
	panic("implement me")
}

func NewRepository(db *sql.DB) service.Repository {
	return &Repository{
		db: db,
	}
}
