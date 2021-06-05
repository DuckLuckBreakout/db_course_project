package repository

import (
	"database/sql"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) forum.Repository {
	return &Repository{
		db: db,
	}
}
