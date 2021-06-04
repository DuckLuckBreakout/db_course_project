package repository

import (
	"database/sql"
	"github.com/DuckLuckBreakout/db_course_project/internal/errors"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/models"
	"github.com/DuckLuckBreakout/db_course_project/internal/pkg/user"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) user.Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Create(user *models.User) error {
	row := r.db.QueryRow(
		"INSERT INTO users(nickname, fullname, about, email) "+
			"VALUES ($1, $2, $3, $4)",
		user.Nickname,
		user.Fullname,
		user.About,
		user.Email,
	)

	var userId uint64
	if err := row.Scan(&userId); err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetAllUsersByNicknameAndEmail(user *models.User) ([]*models.User, error) {
	rows, _ := r.db.Query(
		"SELECT nickname, fullname, about, email "+
			"FROM users "+
			"WHERE nickname = $1 OR email = $2",
		user.Nickname,
		user.Email,
	)

	users := make([]*models.User, 0)
	for rows.Next() {
		rowUser := &models.User{}
		err := rows.Scan(
			&rowUser.Nickname,
			&rowUser.Fullname,
			&rowUser.About,
			&rowUser.Email,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, rowUser)
	}

	return users, nil
}

func (r *Repository) GetUserByNickname(user *models.User) error {
	row := r.db.QueryRow(
		"SELECT fullname, about, email "+
			"FROM users "+
			"WHERE nickname = $1",
		user.Nickname,
	)

	if err := row.Scan(
		&user.Fullname,
		&user.About,
		&user.Email,
	); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(user *models.User) error {
	row, err := r.db.Exec(
		"UPDATE users " +
			"SET fullname = $2, about = $3, email = $4 "+
			"WHERE nickname = $1",
		user.Nickname,
		user.Fullname,
		user.About,
		user.Email,
	)

	if err != nil {
		return err
	}

	affectedRowsCount, _ := row.RowsAffected()
	if affectedRowsCount == 0 {
		return errors.ErrUserNotFound
	}

	return nil
}