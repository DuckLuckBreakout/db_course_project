package repository

import (
	"database/sql"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) thread.Repository {
	return &Repository{
		db: db,
	}
}
