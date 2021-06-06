package repository

import (
	"database/sql"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread"
)

type Repository struct {
	db *sql.DB
}

func (r Repository) Details(thread *models.Thread) error {
	var threadId int32
	if thread.Id == 0 {
		row := r.db.QueryRow("SELECT id "+
			"FROM threads "+
			"WHERE slug = $1", thread.Slug)
		if err := row.Err(); err != nil {
			return err
		}
		if err := row.Scan(&threadId); err != nil {
			return err
		}
	} else {
		threadId = thread.Id
	}

	threadRow := r.db.QueryRow("SELECT author, created, forum, id, message, slug, title, votes "+
		"FROM threads "+
		"WHERE id = $1", threadId)
	if err := threadRow.Err(); err != nil {
		return err
	}

	if err := threadRow.Scan(
		&thread.Author,
		&thread.Created,
		&thread.Forum,
		&thread.Id,
		&thread.Message,
		&thread.Slug,
		&thread.Title,
		&thread.Votes,
	); err != nil {
		return err
	}
	return nil
}

func (r Repository) Vote(thread *models.ThreadVoice) (*models.Thread, error){
	var threadId int32
	if thread.ThreadID == 0 {
		row := r.db.QueryRow("SELECT id "+
			"FROM threads "+
			"WHERE slug = $1", thread.ThreadSlug)
		if err := row.Err(); err != nil {
			return nil, err
		}
		if err := row.Scan(&threadId); err != nil {
			return nil, err
		}
	} else {
		threadId = thread.ThreadID
	}
	row, err := r.db.Exec("UPDATE voices "+
		"SET voice = $1 "+
		"WHERE nickname = $2 AND thread = $3", thread.Voice, thread.Nickname, threadId)
	if err != nil {
		return nil, err
	}
	if affected, _ := row.RowsAffected(); affected == 0 {
		row := r.db.QueryRow("INSERT INTO voices(nickname, voice, thread) "+
			"VALUES ($1, $2, $3)", thread.Nickname, thread.Voice, threadId)
		if err := row.Err(); err != nil {
			return nil, err
		}
	}

	threadRow := r.db.QueryRow("SELECT author, created, forum, id, message, slug, title, votes "+
		"FROM threads "+
		"WHERE id = $1", threadId)
	if err := threadRow.Err(); err != nil {
		return nil, err
	}

	var threadInfo models.Thread
	if err := threadRow.Scan(
		&threadInfo.Author,
		&threadInfo.Created,
		&threadInfo.Forum,
		&threadInfo.Id,
		&threadInfo.Message,
		&threadInfo.Slug,
		&threadInfo.Title,
		&threadInfo.Votes,
	); err != nil {
		return nil, err
	}
	return &threadInfo, nil
}

func NewRepository(db *sql.DB) thread.Repository {
	return &Repository{
		db: db,
	}
}
