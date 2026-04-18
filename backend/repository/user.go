package repository

import (
	"database/sql"
	"tsumiki/schema"
)

type UserRepository interface {
	FindByID(id int) (*schema.User, error)
}

type userRepositoryImpl struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (ur *userRepositoryImpl) FindByID(id int) (*schema.User, error) {
	var user schema.User
	err := ur.db.QueryRow(
		"SELECT id, discord_user_id, name, avatar_url, created_at, updated_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.DiscordUserID, &user.Name, &user.AvatarUrl, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
