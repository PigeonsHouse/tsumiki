package repository

import (
	"database/sql"
	"tsumiki/schema"
)

type AuthRepository interface {
	FindByDiscordUserId(id string) (*schema.User, error)
	CreateUserByDiscord(
		name string,
		thumbnail_url string,
		discord_user_id string,
		guild_id string,
	) (*schema.User, error)
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

	return nil, nil
}

func (ar *authRepositoryImpl) CreateUserByDiscord(
	name string,
	thumbnail_url string,
	discord_user_id string,
	guild_id string,
) (*schema.User, error) {
	return nil, nil
}
