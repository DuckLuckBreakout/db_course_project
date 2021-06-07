package repository

import (
	"database/sql"
	"fmt"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strings"
	"time"
)

type Repository struct {
	db *sql.DB
}

func (r Repository) Threads(thread *models.ThreadSearch) ([]*models.Thread, error) {
	checkForum := r.db.QueryRow("SELECT COUNT(*) "+
		"FROM forums "+
		"WHERE slug = $1", thread.Forum)
	var checkResult int
	if err := checkForum.Scan(&checkResult); err != nil {
		return nil, err
	}
	if checkResult == 0 {
		return nil, errors.ErrUserNotFound
	}

	var sortDirection string
	if thread.Desc {
		since, _ := time.Parse("2006-01-02T15:04:05.000Z", "3006-01-02T15:04:05.000Z")
		if thread.Since != since {
			sortDirection = "AND created <= $2 ORDER BY created DESC "
		} else {
			sortDirection = "ORDER BY created DESC "
		}
	} else {
		since, _ := time.Parse("2006-01-02T15:04:05.000Z", "3006-01-02T15:04:05.000Z")
		if thread.Since != since {
			sortDirection = "AND created >= $2 ORDER BY created ASC "
		} else {
			sortDirection = "ORDER BY created ASC "
		}
	}

	if strings.Contains(sortDirection, "$2") {
		rows, err := r.db.Query("SELECT id, title, author, forum, message, votes, slug, created "+
			"FROM threads "+
			"WHERE forum = $1 "+sortDirection+" "+
			"LIMIT $3", thread.Forum, thread.Since, thread.Limit)
		if err != nil {
			return nil, err
		}
		threads := make([]*models.Thread, 0)
		for rows.Next() {
			rowThread := &models.Thread{}
			err := rows.Scan(
				&rowThread.Id,
				&rowThread.Title,
				&rowThread.Author,
				&rowThread.Forum,
				&rowThread.Message,
				&rowThread.Votes,
				&rowThread.Slug,
				&rowThread.Created,
			)
			if err != nil {
				return nil, err
			}
			threads = append(threads, rowThread)
		}
		return threads, nil
	} else {
		rows, err := r.db.Query("SELECT id, title, author, forum, message, votes, slug, created "+
			"FROM threads "+
			"WHERE forum = $1 "+sortDirection+" "+
			"LIMIT $2", thread.Forum, thread.Limit)
		if err != nil {
			return nil, err
		}
		threads := make([]*models.Thread, 0)
		for rows.Next() {
			rowThread := &models.Thread{}
			err := rows.Scan(
				&rowThread.Id,
				&rowThread.Title,
				&rowThread.Author,
				&rowThread.Forum,
				&rowThread.Message,
				&rowThread.Votes,
				&rowThread.Slug,
				&rowThread.Created,
			)
			if err != nil {
				return nil, err
			}
			threads = append(threads, rowThread)
		}
		return threads, nil
	}
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
			row = r.db.QueryRow("SELECT slug "+
				"FROM forums "+
				"WHERE slug = $1", thread.Forum)
			if err := row.Scan(&thread.Forum); err != nil {
				return err
			}
			return errors.ErrThreadAlreadyCreatedError
		}
	}
	row := r.db.QueryRow("SELECT slug "+
		"FROM forums "+
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
		fmt.Println(err)
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
			fmt.Println(err)
			return errors.ErrUserNotFound
		}
		row := r.db.QueryRow("SELECT title, \"user\", slug, posts, threads "+
			"FROM forums "+
			"WHERE slug = $1", forum.Slug)
		if err := row.Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads); err != nil {
			fmt.Println(err)
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
