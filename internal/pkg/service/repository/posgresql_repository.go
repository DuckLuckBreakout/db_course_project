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
	_, err := r.db.Exec("TRUNCATE TABLE users CASCADE")
	if err != nil {
		return err
	}
	_, err = r.db.Exec("TRUNCATE TABLE posts CASCADE")
	if err != nil {
		return err
	}
	_, err = r.db.Exec("TRUNCATE TABLE threads CASCADE")
	if err != nil {
		return err
	}
	_, err = r.db.Exec("TRUNCATE TABLE forums CASCADE ")
	if err != nil {
		return err
	}
	return nil
}

func (r Repository) Status() (*models.Status, error) {
	var status models.Status
	row := r.db.QueryRow("SELECT COUNT(*) FROM users")
	if err := row.Scan(&status.User); err != nil {
		return nil, err
	}
	row = r.db.QueryRow("SELECT COUNT(*) FROM posts")
	if err := row.Scan(&status.Post); err != nil {
		return nil, err
	}
	row = r.db.QueryRow("SELECT COUNT(*) FROM forums")
	if err := row.Scan(&status.Forum); err != nil {
		return nil, err
	}
	row = r.db.QueryRow("SELECT COUNT(*) FROM threads")
	if err := row.Scan(&status.Thread); err != nil {
		return nil, err
	}
	return &status, nil
}

func NewRepository(db *sql.DB) service.Repository {
	return &Repository{
		db: db,
	}
}
