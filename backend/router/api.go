package router

import (
	"tsumiki/handler"

	"github.com/go-chi/chi/v5"
)

func SetApiRouter(mux *chi.Mux, handlers *handler.Handlers) {
	mux.Get("/api/v1/ping", handlers.Ping.Ping)

	mux.Get("/api/v1/auth/discord", handlers.Auth.RedirectDiscord)
	mux.Get("/api/v1/auth/discord/callback", handlers.Auth.CallbackDiscord)
}
