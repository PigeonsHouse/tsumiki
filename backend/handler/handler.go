package handler

import (
	"tsumiki/media"
	"tsumiki/repository"
	"tsumiki/store"
)

type Handlers struct {
	Ping PingHandler
	Auth AuthHandler
}

func NewHandlers(repos *repository.Repositories, stores *store.Stores, mediaSvc media.MediaService) *Handlers {
	return &Handlers{
		Ping: NewPingHandler(),
		Auth: NewAuthHandler(repos.Auth, stores.Auth, mediaSvc),
	}
}
