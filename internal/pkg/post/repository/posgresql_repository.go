package repository

import (
	"database/sql"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/post"
)

type Repository struct {
	db *sql.DB
}

func (r Repository) DetailsUser(id int) (*models.User, error) {
	row := r.db.QueryRow("SELECT author "+
		"FROM posts "+
		"WHERE id = $1", id)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var author string
	if err := row.Scan(
		&author,
	); err != nil {
		return nil, errors.ErrUserNotFound
	}

	row = r.db.QueryRow("SELECT nickname, fullname, about, email "+
		"FROM users "+
		"WHERE nickname = $1", author)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var userInfo models.User
	if err := row.Scan(
		&userInfo.Nickname,
		&userInfo.Fullname,
		&userInfo.About,
		&userInfo.Email,
	); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

func (r Repository) DetailsForum(id int) (*models.Forum, error) {
	row := r.db.QueryRow("SELECT forum "+
		"FROM posts "+
		"WHERE id = $1", id)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var slug string
	if err := row.Scan(
		&slug,
	); err != nil {
		return nil, errors.ErrUserNotFound
	}

	row = r.db.QueryRow("SELECT  title, \"user\", slug, posts, threads "+
		"FROM forums "+
		"WHERE slug = $1", slug)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var forumInfo models.Forum
	if err := row.Scan(
		&forumInfo.Title,
		&forumInfo.User,
		&forumInfo.Slug,
		&forumInfo.Posts,
		&forumInfo.Threads,
	); err != nil {
		return nil, err
	}

	return &forumInfo, nil
}

func (r Repository) DetailsThread(id int) (*models.Thread, error) {

	row := r.db.QueryRow("SELECT thread "+
		"FROM posts "+
		"WHERE id = $1", id)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var thread int64
	if err := row.Scan(
		&thread,
	); err != nil {
		return nil, errors.ErrUserNotFound
	}

	row = r.db.QueryRow("SELECT  id, title, author, forum, message, votes, slug, created "+
		"FROM threads "+
		"WHERE id = $1", thread)
	if err := row.Err(); err != nil {
		return nil, err
	}

	var threadInfo models.Thread
	if err := row.Scan(
		&threadInfo.Id,
		&threadInfo.Title,
		&threadInfo.Author,
		&threadInfo.Forum,
		&threadInfo.Message,
		&threadInfo.Votes,
		&threadInfo.Slug,
		&threadInfo.Created,
	); err != nil {
		return nil, err
	}

	return &threadInfo, nil
}

func (r Repository) UpdateDetails(updatePost *models.Post) (*models.Post, error) {
	if updatePost.Message == "" {
		row := r.db.QueryRow("SELECT id, parent, author, message, is_edited, forum, thread, created "+
			"FROM posts "+
			"WHERE id = $1", updatePost.Id)
		if err := row.Err(); err != nil {
			return nil, errors.ErrUserNotFound
		}
		if err := row.Scan(
			&updatePost.Id,
			&updatePost.Parent,
			&updatePost.Author,
			&updatePost.Message,
			&updatePost.IsEdited,
			&updatePost.Forum,
			&updatePost.Thread,
			&updatePost.Created,
		); err != nil {
			return nil, errors.ErrUserNotFound
		}
		return nil, nil
	}
	row := r.db.QueryRow("SELECT message "+
		"FROM posts "+
		"WHERE id = $1", updatePost.Id)
	if row.Err() != nil {
		return nil, errors.ErrUserNotFound
	}

	var message string
	if err := row.Scan(&message); err != nil {
		return nil, errors.ErrUserNotFound
	}
	if message != updatePost.Message {
		_, err := r.db.Exec("UPDATE posts "+
			"SET is_edited = $1, message=$2 "+
			"WHERE id = $3", true, updatePost.Message, updatePost.Id)
		if err != nil {
			return nil, errors.ErrUserNotFound
		}
	}

	row = r.db.QueryRow("SELECT id, parent, author, message, is_edited, forum, thread, created "+
		"FROM posts "+
		"WHERE id = $1", updatePost.Id)
	if err := row.Err(); err != nil {
		return nil, errors.ErrUserNotFound
	}
	if err := row.Scan(
		&updatePost.Id,
		&updatePost.Parent,
		&updatePost.Author,
		&updatePost.Message,
		&updatePost.IsEdited,
		&updatePost.Forum,
		&updatePost.Thread,
		&updatePost.Created,
	); err != nil {
		return nil, errors.ErrUserNotFound
	}
	//if message == updatePost.Message {
	//	updatePost.IsEdited = false
	//}
	return nil, nil
}

func (r Repository) Details(id int) (*models.Post, error) {
	row := r.db.QueryRow("SELECT author, created, forum, id, message, thread, is_edited "+
		"FROM posts "+
		"WHERE id = $1", id)
	if err := row.Err(); err != nil {
		return nil, errors.ErrUserNotFound
	}
	var postInfo models.Post
	if err := row.Scan(
		&postInfo.Author,
		&postInfo.Created,
		&postInfo.Forum,
		&postInfo.Id,
		&postInfo.Message,
		&postInfo.Thread,
		&postInfo.IsEdited,
	); err != nil {
		return nil, errors.ErrUserNotFound
	}
	return &postInfo, nil
}

func NewRepository(db *sql.DB) post.Repository {
	return &Repository{
		db: db,
	}
}
