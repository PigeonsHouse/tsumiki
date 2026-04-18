package main

import (
	"context"
	"fmt"
	"net/http"
	"tsumiki/env"
	"tsumiki/external"
	"tsumiki/handler"
	"tsumiki/infra"
	"tsumiki/media"
	"tsumiki/repository"
	"tsumiki/router"
	"tsumiki/store"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := env.LoadEnv(); err != nil {
		panic(fmt.Errorf("env: %w", err))
	}

	db, err := infra.NewDatabase()
	if err != nil {
		panic(fmt.Errorf("db: %w", err))
	}
	redis, err := infra.NewRedis()
	if err != nil {
		panic(fmt.Errorf("redis: %w", err))
	}
	s3Client, err := infra.NewS3Client()
	if err != nil {
		panic(fmt.Errorf("s3: %w", err))
	}
	mediaSvc, err := media.NewMediaService(context.Background(), s3Client)
	if err != nil {
		panic(fmt.Errorf("media: %w", err))
	}
	stores := store.NewStores(redis)
	repos := repository.NewRepositories(db)
	discordSvc := external.NewDiscordService()
	handlers := handler.NewHandlers(repos, stores, mediaSvc, discordSvc)

	mux := chi.NewRouter()
	router.SetApiRouter(mux, handlers)
	router.SetFrontendRouter(mux, "./view")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", env.AppPort), mux); err != nil {
		panic(err)
	}
}
