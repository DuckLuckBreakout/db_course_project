package repository

import (
	"database/sql"
	"fmt"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/service"
)

type Repository struct {
	db *sql.DB
}

func (r *Repository) Close() {
	row := r.db.QueryRow("SELECT pg_terminate_backend(pid) FROM pg_stat_activity " +
		"WHERE datname = 'forum' " +
		"AND pid <> pg_backend_pid() " +
		"AND state in ('idle')")
	if row.Err() != nil {
		fmt.Println(row.Err())
	}
}

func (r Repository) Clear() error {
	row := r.db.QueryRow("TRUNCATE TABLE users CASCADE")
	if err := row.Err(); err != nil {
		return err
	}
	row = r.db.QueryRow("TRUNCATE TABLE posts CASCADE")
	if err := row.Err(); err != nil {
		return err
	}
	row = r.db.QueryRow("TRUNCATE TABLE threads CASCADE")
	if err := row.Err(); err != nil {
		return err
	}
	row = r.db.QueryRow("TRUNCATE TABLE forums CASCADE ")
	if err := row.Err(); err != nil {
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
