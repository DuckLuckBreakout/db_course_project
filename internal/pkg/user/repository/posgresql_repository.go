package repository

import (
	"database/sql"
	"fmt"
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
	_, err := r.db.Exec(
		"INSERT INTO users(nickname, fullname, about, email) "+
			"VALUES ($1, $2, $3, $4)",
		user.Nickname,
		user.Fullname,
		user.About,
		user.Email,
	)

	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *Repository) GetAllUsersByNicknameAndEmail(user *models.User) ([]*models.User, error) {
	rows, err := r.db.Query(
		"SELECT nickname, fullname, about, email "+
			"FROM users "+
			"WHERE nickname = $1 OR email = $2",
		user.Nickname,
		user.Email,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
			fmt.Println(err)
			return nil, err
		}
		users = append(users, rowUser)
	}

	return users, nil
}

func (r *Repository) GetUserByNickname(user *models.User) error {

	row := r.db.QueryRow(
		"SELECT nickname, fullname, about, email "+
			"FROM users "+
			"WHERE nickname = $1",
		user.Nickname,
	)

	if err := row.Scan(
		&user.Nickname,
		&user.Fullname,
		&user.About,
		&user.Email,
	); err != nil {
		return err
	}

	return nil
}

func (r *Repository) Update(user *models.User) error {

	setString := "SET "
	if user.Fullname != "" {
		setString += " fullname = " + "'" + user.Fullname + "' "
		if user.About != "" || user.Email != "" {
			setString += ", "
		}
	}
	if user.About != "" {
		setString += " about = " + "'" + user.About + "' "
		if user.Email != "" {
			setString += ", "
		}
	}
	if user.Email != "" {
		setString += " email = " + "'" + user.Email + "' "
	}

	row, err := r.db.Exec(
		"UPDATE users "+
			setString+
			"WHERE nickname = $1",
		user.Nickname,
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
