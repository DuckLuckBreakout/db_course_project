package repository

import (
	"database/sql"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

type Repository struct {
	db *sql.DB
}

func (r Repository) Create(forum *models.Forum) error {
	row := r.db.QueryRow("SELECT nickname "+
		"FROM users "+
		"WHERE nickname = $1", forum.User)
	if err := row.Scan(&forum.User); err != nil {
		return errors.ErrUserNotFound
	}

	row = r.db.QueryRow(
		"INSERT INTO forums(title, \"user\", slug) "+
			"VALUES ($1, $2, $3)",
		forum.Title,
		forum.User,
		forum.Slug,
	)
	err := row.Err()
	if err != nil {
		if err.(*pq.Error).Code == "23503" {
			return errors.ErrUserNotFound
		}
		row := r.db.QueryRow("SELECT title, \"user\", slug, posts, threads "+
			"FROM forums "+
			"WHERE slug = $1", forum.Slug)
		if err := row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads); err != nil {
			return err
		}
		return errors.ErrForumAlreadyCreatedError
	}
	return nil
}

func NewRepository(db *sql.DB) forum.Repository {
	return &Repository{
		db: db,
	}
}
