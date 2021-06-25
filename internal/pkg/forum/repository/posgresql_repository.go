package repository

import (
	"database/sql"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/forum"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/bradfitz/slice"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strings"
)

type Repository struct {
	db *sql.DB
}

func (r Repository) Users(searchParams *models.UserSearch) ([]*models.User, error) {
	row := r.db.QueryRow("SELECT slug "+
		"FROM forums "+
		"WHERE slug = $1", searchParams.Forum)
	if err := row.Err(); err != nil {
		return nil, err
	}
	var forumSlug string
	if err := row.Scan(&forumSlug); err != nil {
		return nil, err
	}

	var sortChar string
	if searchParams.Desc {
		//sortString = " DESC "
		sortChar = " < "
		if searchParams.Since == "" {
			searchParams.Since = "ZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZZ"
		}
	} else {
		//sortString = " ASC "
		sortChar = " > "
	}

	rows, err := r.db.Query("SELECT author "+
		"FROM posts "+
		"WHERE forum = $1 AND author "+sortChar+" $2 ", searchParams.Forum, searchParams.Since)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make(map[string]struct{}, 0)
	var void struct{}

	for rows.Next() {
		var user string
		err := rows.Scan(&user)
		if err != nil {
			return nil, err
		}
		users[user] = void
	}

	rows1, err := r.db.Query("SELECT author "+
		"FROM threads "+
		"WHERE forum = $1 AND author "+sortChar+" $2 ", searchParams.Forum, searchParams.Since)
	if err != nil {
		return nil, err
	}
	defer rows1.Close()

	for rows1.Next() {
		var user string
		err := rows1.Scan(&user)
		if err != nil {
			return nil, err
		}
		users[user] = void
	}
	nickNames := make([]string, 0)
	for user, _ := range users {
		nickNames = append(nickNames, user)
	}

	if sortChar == " > " {
		slice.Sort(nickNames[:], func(i, j int) bool {
			return strings.ToLower(nickNames[i]) < strings.ToLower(nickNames[j])
		})
	} else {
		slice.Sort(nickNames[:], func(i, j int) bool {
			return strings.ToLower(nickNames[i]) > strings.ToLower(nickNames[j])
		})
	}

	if searchParams.Limit == 0 {
		searchParams.Limit = 100
	}

	sliceLen := int32(len(nickNames))
	if searchParams.Limit > sliceLen {
		searchParams.Limit = sliceLen
	}

	nickNames = nickNames[:searchParams.Limit]

	usersInfo := make([]*models.User, 0)
	for _, user := range nickNames {
		row := r.db.QueryRow("SELECT about, email, fullname, nickname "+
			"FROM users "+
			"WHERE nickname = $1", user)
		if err := row.Err(); err != nil {
			return nil, err
		}
		var userInfo models.User
		if err := row.Scan(&userInfo.About, &userInfo.Email, &userInfo.Fullname, &userInfo.Nickname); err != nil {
			return nil, err
		}
		usersInfo = append(usersInfo, &userInfo)
	}

	return usersInfo, nil
}

func (r Repository) Threads(thread *models.ThreadSearch, sinceString string) ([]*models.Thread, error) {
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
		//since, _ := time.Parse("2006-01-02T15:04:05.000Z", "3006-01-02T15:04:05.000Z")
		if sinceString != "" {
			sortDirection = "AND created <= $2 ORDER BY created DESC "
		} else {
			sortDirection = "ORDER BY created DESC "
		}
	} else {
		//since, _ := time.Parse("2006-01-02T15:04:05.000Z", "3006-01-02T15:04:05.000Z")
		if sinceString != "" {
			sortDirection = "AND created >= $2 ORDER BY created ASC "
		} else {
			sortDirection = "ORDER BY created ASC "
		}
	}

	if strings.Contains(sortDirection, "$2") {
		rows, err := r.db.Query("SELECT id, title, author, forum, message, votes, slug, created "+
			"FROM threads "+
			"WHERE forum = $1 "+sortDirection+" "+
			"LIMIT $3", thread.Forum, sinceString, thread.Limit)
		if err != nil {
			if rows != nil {
				rows.Close()
			}
			return nil, err
		}
		defer rows.Close()

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
			if rows != nil {
				rows.Close()
			}
			return nil, err
		}
		defer rows.Close()

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
		"VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, created", thread.Title, thread.Author, thread.Message, thread.Forum, thread.Slug, thread.Created)
	if err := row.Err(); err != nil {
		return err
	}

	if err := row.Scan(&thread.Id, &thread.Created); err != nil {
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

	_, err := r.db.Exec(
		"INSERT INTO forums(title, \"user\", slug) "+
			"VALUES ($1, $2, $3)",
		forum.Title,
		forum.User,
		forum.Slug,
	)
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
