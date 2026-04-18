package router

import (
	"tsumiki/handler"
	"tsumiki/middleware"

	"github.com/go-chi/chi/v5"
)

func SetApiRouter(mux *chi.Mux, handlers *handler.Handlers) {
	mux.Get("/api/v1/ping", handlers.Ping.Ping)

	mux.Get("/api/v1/auth/discord", handlers.Auth.RedirectDiscord)
	mux.Get("/api/v1/auth/discord/callback", handlers.Auth.CallbackDiscord)
	mux.Get("/api/v1/auth/token-refresh", handlers.Auth.RefreshToken)

	mux.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/api/v1/users/me", handlers.User.GetMyInfo)
		r.Get("/api/v1/users/{userID}", handlers.User.GetUserInfo)
	})
}
