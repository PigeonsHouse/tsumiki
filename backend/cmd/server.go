package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"tsumiki/env"
	"tsumiki/handler"
	"tsumiki/repository"
	"tsumiki/router"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := env.LoadEnv(); err != nil {
		panic(err)
	}

	// TODO: infraパッケージとかを切る
	db, _ := sql.Open("mysql", "")
	repos := repository.NewRepositories(db)
	handlers := handler.NewHandlers(repos)

	mux := chi.NewRouter()
	router.SetApiRouter(mux, handlers)
	router.SetFrontendRouter(mux, "./view")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", env.AppPort), mux); err != nil {
		panic(err)
	}
}
