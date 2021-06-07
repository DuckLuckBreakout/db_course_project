package repository

import (
	"database/sql"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread"
	"strconv"
)

type Repository struct {
	db *sql.DB
}

func (r Repository) Create(slugOrId string, posts []*models.Post) error {
	var threadId int32
	var forum string
	id, err := strconv.Atoi(slugOrId)
	if err != nil {
		row := r.db.QueryRow("SELECT id, forum "+
			"FROM threads "+
			"WHERE slug = $1", slugOrId)
		if err := row.Err(); err != nil {
			return err
		}
		if err := row.Scan(&threadId, &forum); err != nil {
			return err
		}
	} else {
		threadId = int32(id)
		row := r.db.QueryRow("SELECT forum "+
			"FROM threads "+
			"WHERE id = $1", threadId)
		if err := row.Err(); err != nil {
			return err
		}
		if err := row.Scan(&forum); err != nil {
			return err
		}
	}

	if len(posts) == 0 {
		return nil
	}

	query := "INSERT INTO posts(parent, author, message, thread, forum) " +
		"VALUES "
	varIndex := 1
	var values []interface{}
	for i, post := range posts {
		if i != 0 {
			query += ", "
		}
		post.Thread = threadId
		post.Forum = forum
		query += "($" + strconv.Itoa(varIndex) + ", $" + strconv.Itoa(varIndex+1) +
			", $" + strconv.Itoa(varIndex+2) + ", $" + strconv.Itoa(varIndex+3) + ", $" + strconv.Itoa(varIndex+4) + ") "
		varIndex += 5
		values = append(values, post.Parent, post.Author, post.Message, post.Thread, post.Forum)
		if post.Parent != 0 {
			row := r.db.QueryRow("SELECT COUNT(*) " +
				"FROM posts " +
				"WHERE id = $1", post.Parent)
			var parentId int64
			if err := row.Scan(&parentId); err != nil {
				return err
			}
			if parentId == 0 {
				return errors.ErrUserNotFound
			}
		}

	}
	rows, err := r.db.Query(query + " RETURNING id", values...)
	if err != nil {
		return err
	}
	i := 0
	for rows.Next() {
		err := rows.Scan(
			&posts[i].Id,
		)
		i++
		if err != nil {
			return err
		}
	}
	return nil
}

func (r Repository) UpdateDetails(thread *models.ThreadUpdate) (*models.Thread, error) {
	var threadId int32
	if thread.Id == 0 {
		row := r.db.QueryRow("SELECT id "+
			"FROM threads "+
			"WHERE slug = $1", thread.Slug)
		if err := row.Err(); err != nil {
			return nil, err
		}
		if err := row.Scan(&threadId); err != nil {
			return nil, err
		}
	} else {
		threadId = thread.Id
	}

	row, err := r.db.Exec("UPDATE threads "+
		"SET title = $1, message = $2 "+
		"WHERE id = $3", thread.Title, thread.Message, threadId)
	if err != nil {
		return nil, err
	}
	if affected, _ := row.RowsAffected(); affected == 0 {
		return nil, err
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

func (r Repository) Vote(thread *models.ThreadVoice) (*models.Thread, error) {
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
