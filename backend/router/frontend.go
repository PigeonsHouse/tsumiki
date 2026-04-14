package router

import (
	"net/http"
	"path"
	"path/filepath"

	"github.com/go-chi/chi/v5"
)

func SetFrontendRouter(mux *chi.Mux, filePath string) {
	fileServer := http.FileServer(http.Dir(filePath))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, file := path.Split(r.URL.Path)
		ext := filepath.Ext(file)
		// SPA フォールバック
		if file == "" || ext == "" {
			r.URL.Path = "/"
		}
		fileServer.ServeHTTP(w, r)
	})
}
