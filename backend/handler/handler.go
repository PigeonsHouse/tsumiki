package handler

import (
	"tsumiki/external"
	"tsumiki/media"
	"tsumiki/repository"
	"tsumiki/store"
)

type Handlers struct {
	Ping    PingHandler
	Auth    AuthHandler
	User    UserHandler
	Tsumiki TsumikiHandler
	Work    WorkHandler
}

func NewHandlers(
	repos *repository.Repositories,
	stores *store.Stores,
	mediaSvc media.MediaService,
	discordSvc external.DiscordService,
) *Handlers {
	return &Handlers{
		Ping:    NewPingHandler(),
		Auth:    NewAuthHandler(repos, stores.Auth, mediaSvc, discordSvc),
		User:    NewUserHandler(repos, mediaSvc),
		Tsumiki: NewTsumikiHandler(repos, mediaSvc),
		Work:    NewWorkHandler(repos),
	}
}
