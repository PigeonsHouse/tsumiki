package router

import (
	"tsumiki/handler"
	"tsumiki/middleware"

	"github.com/go-chi/chi/v5"
)

func SetApiRouter(mux *chi.Mux, handlers *handler.Handlers) {
	mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", handlers.Ping.Ping)

		r.Route("/auth", func(r chi.Router) {
			r.Get("/discord", handlers.Auth.RedirectDiscord)
			r.Get("/discord/callback", handlers.Auth.CallbackDiscord)
			r.Get("/token-refresh", handlers.Auth.RefreshToken)
		})

		r.Route("/users", func(r chi.Router) {
			r.Use(middleware.AuthMiddleware)
			r.Get("/me", handlers.User.GetMyInfo)
			r.Get("/{userID}", handlers.User.GetUserInfo)
		})
	})
}
