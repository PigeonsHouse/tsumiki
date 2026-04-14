package router

import (
	"tsumiki/handler"

	"github.com/go-chi/chi/v5"
)

func SetApiRouter(mux *chi.Mux) {
	mux.Get("/api/v1/ping", handler.Ping)

	mux.Get("/api/v1/auth/discord", handler.RedirectDiscord)
}
