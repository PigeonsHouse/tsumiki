package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")

	mux := http.NewServeMux()
	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]any{
			"status": "ok",
		})
	})

	fileServer := http.FileServer(http.Dir("./view"))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, file := path.Split(r.URL.Path)
		ext := filepath.Ext(file)
		// SPA フォールバック
		if file == "" || ext == "" {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		panic(err)
	}
}
