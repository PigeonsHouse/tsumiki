package handler

import "tsumiki/repository"

type Handlers struct {
	Ping PingHandler
	Auth AuthHandler
}

func NewHandlers(repos *repository.Repositories) *Handlers {
	return &Handlers{
		Ping: NewPingHandler(),
		Auth: NewAuthHandler(repos.Auth),
	}
}
