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

func (r Repository) CreateThread(thread *models.Thread) error {
	if thread.Slug != "" {
		row := r.db.QueryRow("SELECT COUNT(*) "+
			"FROM threads "+
			"WHERE slug = $1", thread.Slug)
		var result int
		if err := row.Scan(&result); err != nil {
			return err
		}
		if result != 0 {
			row = r.db.QueryRow("SELECT id, title, author, forum, message, votes, slug, created "+
				"FROM threads "+
				"WHERE slug = $1", thread.Slug)
			if err := row.Scan(
				&thread.Id,
				&thread.Title,
				&thread.Author,
				&thread.Forum,
				&thread.Message,
				&thread.Votes,
				&thread.Slug,
				&thread.Created,
			); err != nil {
				return err
			}
			row = r.db.QueryRow("SELECT slug " +
				"FROM forums " +
				"WHERE slug = $1", thread.Forum)
			if err := row.Scan(&thread.Forum); err != nil {
				return err
			}
			return errors.ErrThreadAlreadyCreatedError
		}
	}
	row := r.db.QueryRow("SELECT slug " +
		"FROM forums " +
		"WHERE slug = $1", thread.Forum)
	if err := row.Scan(&thread.Forum); err != nil {
		return err
	}
	row = r.db.QueryRow("INSERT INTO threads(title, author, message, forum, slug, created) "+
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id", thread.Title, thread.Author, thread.Message, thread.Forum, thread.Slug, thread.Created)
	if err := row.Err(); err != nil {
		return err
	}

	if err := row.Scan(&thread.Id); err != nil {
		return err
	}
	return nil
}

func (r Repository) Details(forum *models.Forum) error {
	row := r.db.QueryRow("SELECT title, \"user\", slug, posts, threads "+
		"FROM forums "+
		"WHERE slug = $1", forum.Slug)
	if err := row.Err(); err != nil {
		return err
	}
	if err := row.Scan(
		&forum.Title,
		&forum.User,
		&forum.Slug,
		&forum.Posts,
		&forum.Threads,
	); err != nil {
		return err
	}
	return nil
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
