package repository

import (
	"database/sql"
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
		"SELECT nickname, fullname, about, email " +
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