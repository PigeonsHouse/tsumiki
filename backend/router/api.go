package router

import (
	"tsumiki/handler"
	"tsumiki/middleware"

	"github.com/go-chi/chi/v5"
)

func SetApiRouter(mux *chi.Mux, handlers *handler.Handlers) {
	// TODO: レートリミットをかける
	mux.Route("/api/v1", func(r chi.Router) {
		r.Get("/ping", handlers.Ping.Ping)

		r.Route("/auth", func(r chi.Router) {
			r.Get("/discord", handlers.Auth.RedirectDiscord)
			r.Get("/discord/callback", handlers.Auth.CallbackDiscord)
			r.Get("/token-refresh", handlers.Auth.RefreshToken)
		})
		r.Route("/users", func(r chi.Router) {
			r.Route("/me", func(r chi.Router) {
				r.Use(middleware.RequireAuth)
				r.Get("/", handlers.User.GetMyInfo)
				r.Get("/tsumikis", handlers.Tsumiki.GetMyTsumikis)
			})
			r.Route("/{userID}", func(r chi.Router) {
				r.Get("/", handlers.User.GetUserInfo)
				r.With(middleware.OptionalAuth).
					Get("/tsumikis", handlers.Tsumiki.GetUserTsumikis)
			})
		})

		// TODO: メディア投稿系は強めのレートリミットをかける
		r.Route("/thumbnails", func(r chi.Router) {
			r.Use(middleware.RequireAuth)
			r.Post("/", handlers.Thumbnail.PostThumbnail)
		})

		r.Route("/tsumikis", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(middleware.OptionalAuth)
				r.Get("/", handlers.Tsumiki.GetTsumikis)
				r.Get("/{tsumikiID}", handlers.Tsumiki.GetSpecifiedTsumiki)
			})
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAuth)
				r.Post("/", handlers.Tsumiki.CreateTsumiki)
				r.Put("/{tsumikiID}", handlers.Tsumiki.EditTsumiki)
				r.Put("/{tsumikiID}/thumbnail", handlers.Tsumiki.UpdateTsumikiThumbnail)
				r.Delete("/{tsumikiID}", handlers.Tsumiki.DeleteTsumiki)
				r.Post("/{tsumikiID}/medias", handlers.Tsumiki.PostMedia)
			})
			r.Route("/{tsumikiID}/blocks", func(r chi.Router) {
				r.With(middleware.OptionalAuth).Get("/", handlers.Tsumiki.GetBlocks)
				r.Group(func(r chi.Router) {
					r.Use(middleware.RequireAuth)
					r.Post("/", handlers.Tsumiki.AddBlock)
					r.Put("/{blockID}", handlers.Tsumiki.EditBlock)
					r.Delete("/{blockID}", handlers.Tsumiki.OmitBlock)
				})
			})
		})

		r.Route("/works", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(middleware.OptionalAuth)
				r.Get("/", handlers.Work.GetWorks)
				r.Get("/{workID}", handlers.Work.GetSpecifiedWork)
				r.Get("/{workID}/tsumikis", handlers.Work.GetWorkTsumiki)
			})
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireAuth)
				r.Post("/", handlers.Work.CreateWork)
				r.Put("/{workID}", handlers.Work.EditWork)
				r.Put("/{workID}/thumbnail", handlers.Work.UpdateWorkThumbnail)
				r.Delete("/{workID}", handlers.Work.DeleteWork)
			})
		})
	})
}
