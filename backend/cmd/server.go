package main

import (
	"fmt"
	"net/http"
	"tsumiki/env"
	"tsumiki/router"

	"github.com/go-chi/chi/v5"
)

func main() {
	if err := env.LoadEnv(); err != nil {
		panic(err)
	}

	mux := chi.NewRouter()
	router.SetApiRouter(mux)
	router.SetFrontendRouter(mux, "./view")

	if err := http.ListenAndServe(fmt.Sprintf(":%d", env.AppPort), mux); err != nil {
		panic(err)
	}
}
