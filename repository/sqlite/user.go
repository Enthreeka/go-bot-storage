package sqlite

import (
	"database/sql"
	"github.com/Enthreeka/go-bot-storage/bot/model"
	"github.com/Enthreeka/go-bot-storage/repository"
)

type userRepository struct {
	*SQLite
}

func UserRepository(SQLite *SQLite) repository.User {
	return &userRepository{
		SQLite,
	}
}

func (u *userRepository) Create(user *model.User) error {
	query := `INSERT INTO "user"(id, nickname,first_name,last_name) VALUES ($1,$2,$3,$4)`

	_, err := u.db.Exec(query, user.ID, user.Nickname, user.FirstName, user.LastName)
	return err
}

func (u *userRepository) GetByID(id int64) (*model.User, error) {
	query := `SELECT id, nickname,first_name,last_name,role FROM "user" WHERE id = $1`
	user := &model.User{}

	err := u.db.QueryRow(query, id).Scan(&user.ID,
		&user.Nickname,
		&user.FirstName,
		&user.LastName,
		&user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, sql.ErrNoRows
		}
		return nil, err
	}

	return user, nil
}
