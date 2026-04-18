package repository

import (
	"database/sql"
	"time"
	"tsumiki/schema"
)

type AuthRepository interface {
	FindByDiscordUserId(id string) (*schema.User, error)
	CreateUserByDiscord(
		name string,
		avatar_url string,
		discord_user_id string,
		guild_id string,
	) (*schema.User, error)
	UpdateAvatarUrl(userID int, avatarUrl string) error
}

type authRepositoryImpl struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) AuthRepository {
	return &authRepositoryImpl{
		db: db,
	}
}

func (ar *authRepositoryImpl) FindByDiscordUserId(id string) (*schema.User, error) {
	var user schema.User

	err := ar.db.QueryRow(
		"SELECT id, discord_user_id, name, avatar_url, created_at, updated_at FROM users WHERE discord_user_id = ?",
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

func (ar *authRepositoryImpl) UpdateAvatarUrl(userID int, avatarUrl string) error {
	_, err := ar.db.Exec(
		"UPDATE users SET avatar_url = ? WHERE id = ?",
		avatarUrl, userID,
	)
	return err
}

func (ar *authRepositoryImpl) CreateUserByDiscord(
	name string,
	avatar_url string,
	discord_user_id string,
	guild_id string,
) (*schema.User, error) {
	result, err := ar.db.Exec(
		"INSERT INTO users (discord_user_id, name, avatar_url, guild_id) VALUES (?, ?, ?, ?)",
		discord_user_id, name, avatar_url, guild_id,
	)
	if err != nil {
		return nil, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &schema.User{
		ID:            int(insertedID),
		DiscordUserID: discord_user_id,
		Name:          name,
		AvatarUrl:  avatar_url,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}
