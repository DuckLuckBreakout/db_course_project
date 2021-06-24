package repository

import (
	"database/sql"
	"fmt"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/thread"
	"strconv"
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

func (r Repository) Posts(thread *models.PostSearch) ([]*models.Post, error) {

	var forum string
	if thread.Thread == 0 {
		row := r.db.QueryRow("SELECT id, forum "+
			"FROM threads "+
			"WHERE slug = $1", thread.ThreadSlug)
		if err := row.Err(); err != nil {
			return nil, err
		}
		if err := row.Scan(&thread.Thread, &forum); err != nil {
			return nil, err
		}
	} else {
		row := r.db.QueryRow("SELECT forum "+
			"FROM threads "+
			"WHERE id = $1", thread.Thread)
		if err := row.Err(); err != nil {
			return nil, err
		}
		if err := row.Scan(&forum); err != nil {
			return nil, err
		}
	}

	var selectString, fromString, sortString, sortChar, limitString string
	var values []interface{}

	if thread.Sort == "" {
		thread.Sort = "flat"
	}
	switch thread.Sort {
	case "flat":
		selectString = "SELECT id, parent, author, message, is_edited, forum, thread, created "

		sortString = "ORDER BY id "
		if thread.Desc {
			sortString += "DESC "
			sortChar = " < "
		} else {
			sortString += "ASC "
			sortChar = " > "
		}

		if thread.Since != 0 {
			fromString = "FROM posts " +
				"WHERE thread = $1 AND id " + sortChar + " $2 "
			limitString = "LIMIT $3 "
			values = append(values, thread.Thread, thread.Since, thread.Limit)
		} else {
			fromString = "FROM posts WHERE thread = $1"
			limitString = "LIMIT $2 "
			values = append(values, thread.Thread, thread.Limit)
		}

	case "tree":
		selectString = "SELECT p1.id, p1.parent, p1.author, p1.message, p1.is_edited, p1.forum, p1.thread, p1.created "

		sortString = "ORDER BY p1.path[1] "
		if thread.Desc {
			sortString += "DESC, "
			sortChar = " < "
		} else {
			sortString += "ASC, "
			sortChar = " > "
		}
		sortString += "p1.path "
		if thread.Desc {
			sortString += "DESC "
		} else {
			sortString += "ASC "
		}

		if thread.Since != 0 {
			fromString = "FROM posts p1 JOIN posts p2 ON (p2.id = $1) WHERE (p1.thread = $2 AND p1.path" + sortChar + "p2.path) "
			limitString = "LIMIT $3 "
			values = append(values, thread.Since, thread.Thread, thread.Limit)
		} else {
			fromString = "FROM posts p1 WHERE (p1.thread = $1) "
			limitString = "LIMIT $2 "
			values = append(values, thread.Thread, thread.Limit)
		}

	case "parent_tree":
		selectString = "SELECT p1.id, p1.parent, p1.author, p1.message, p1.is_edited, p1.forum, p1.thread, p1.created "

		sortString = "ORDER BY p1.path[1] "
		if thread.Desc {
			sortString += "DESC, "
			sortChar = " < "
		} else {
			sortString += "ASC, "
			sortChar = " > "
		}
		sortString += "p1.path "

		if thread.Since != 0 {
			if thread.Desc {
				fromString = "FROM posts p1 WHERE p1.path[1] IN ( " +
					"SELECT id " +
					"FROM posts " +
					"WHERE (thread = $1 AND parent = 0 AND " +
					"path[1] " + sortChar + " (SELECT path[1] FROM posts WHERE id = $2)) ORDER BY id DESC LIMIT $3 ) "
			} else {
				fromString = "FROM posts p1 WHERE p1.path[1] IN ( " +
					"SELECT id " +
					"FROM posts " +
					"WHERE (thread = $1 AND parent = 0 AND " +
					"path[1] " + sortChar + " (SELECT path[1] FROM posts WHERE id = $2)) ORDER BY id ASC LIMIT $3 ) "
			}

			values = append(values, thread.Thread, thread.Since, thread.Limit)
		} else {
			if thread.Desc {
				fromString = "FROM posts p1 WHERE p1.path[1] IN ( SELECT id FROM posts WHERE (thread = $1 AND parent = 0) ORDER BY id " + "DESC" + " LIMIT $2 ) "
			} else {
				fromString = "FROM posts p1 WHERE p1.path[1] IN ( SELECT id FROM posts WHERE (thread = $1 AND parent = 0) ORDER BY id " + "ASC" + " LIMIT $2 ) "
			}
			values = append(values, thread.Thread, thread.Limit)
		}

	default:
		return nil, errors.ErrUserNotFound

	}

	rows, err := r.db.Query(selectString+
		fromString+
		sortString+
		limitString, values...)
	if err != nil {
		if rows != nil {
			rows.Close()
		}
		return nil, err
	}
	defer rows.Close()

	posts := make([]*models.Post, 0)
	for rows.Next() {
		var postFormDb models.Post
		if err := rows.Scan(
			&postFormDb.Id,
			&postFormDb.Parent,
			&postFormDb.Author,
			&postFormDb.Message,
			&postFormDb.IsEdited,
			&postFormDb.Forum,
			&postFormDb.Thread,
			&postFormDb.Created,
		); err != nil {
			return nil, err
		}
		posts = append(posts, &postFormDb)
	}

	return posts, nil
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
			return errors.ErrUserNotFound
		}
		if err := row.Scan(&threadId, &forum); err != nil {
			return errors.ErrUserNotFound
		}
	} else {
		threadId = int32(id)
		row := r.db.QueryRow("SELECT forum "+
			"FROM threads "+
			"WHERE id = $1", threadId)
		if err := row.Err(); err != nil {
			return errors.ErrUserNotFound
		}
		if err := row.Scan(&forum); err != nil {
			return errors.ErrUserNotFound
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
			row := r.db.QueryRow("SELECT COUNT(*) "+
				"FROM posts "+
				"WHERE id = $1", post.Parent)
			var parentId int64
			if err := row.Scan(&parentId); err != nil {
				return errors.ErrUserNotFound
			}
			if parentId == 0 {
				return errors.ErrUserAlreadyCreatedError
			}
		}

	}
	rows, err := r.db.Query(query+" RETURNING id", values...)
	if err != nil {
		return err
	}
	defer rows.Close()

	i := 0
	for rows.Next() {
		err := rows.Scan(
			&posts[i].Id,
		)
		i++
		if err != nil {
			return errors.ErrUserNotFound
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
		if threadId == 0 {
			return nil, errors.ErrUserNotFound
		}
	} else {
		threadId = thread.Id
	}

	var row sql.Result
	var err error

	if thread.Message != "" && thread.Title != "" {
		row, err = r.db.Exec("UPDATE threads "+
			"SET title = $1, message = $2 "+
			"WHERE id = $3", thread.Title, thread.Message, threadId)
		if err != nil {
			return nil, err
		}
	}

	if thread.Message == "" && thread.Title != "" {
		row, err = r.db.Exec("UPDATE threads "+
			"SET title = $1 "+
			"WHERE id = $2", thread.Title, threadId)
		if err != nil {
			return nil, err
		}
	}

	if thread.Message != "" && thread.Title == "" {
		row, err = r.db.Exec("UPDATE threads "+
			"SET message = $1 "+
			"WHERE id = $2", thread.Message, threadId)
		if err != nil {
			return nil, err
		}
	}

	if !(thread.Message == "" && thread.Title == "") {
		if affected, _ := row.RowsAffected(); affected == 0 {
			return nil, errors.ErrUserNotFound
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
