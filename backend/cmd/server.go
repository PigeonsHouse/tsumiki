package main

import (
	"fmt"
	"net/http"
	"tsumiki/env"
	"tsumiki/external"
	"tsumiki/handler"
	"tsumiki/repository"
	"tsumiki/router"
	"tsumiki/store"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := env.LoadEnv(); err != nil {
		panic(fmt.Errorf("env: %w", err))
	}

	db, err := external.NewDatabase()
	if err != nil {
		panic(fmt.Errorf("db: %w", err))
	}
	redis, err := external.NewRedis()
	if err != nil {
		panic(fmt.Errorf("redis: %w", err))
	}
	stores := store.NewStores(redis)
	repos := repository.NewRepositories(db)
	handlers := handler.NewHandlers(repos, stores)

	mux := chi.NewRouter()
	router.SetApiRouter(mux, handlers)
	router.SetFrontendRouter(mux, "./view")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", env.AppPort), mux); err != nil {
		panic(err)
	}
}
