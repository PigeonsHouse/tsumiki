package router

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"tsumiki/env"

	"github.com/go-chi/chi/v5"
)

func SetMediaRouter(mux *chi.Mux) {
	mux.Get("/medias/*", func(w http.ResponseWriter, r *http.Request) {
		// /medias/* の * 部分を取得してS3パスに変換
		path := chi.URLParam(r, "*")

		targetURL, err := url.Parse(fmt.Sprintf("%s/%s/%s", env.S3Endpoint, env.S3Bucket, path))
		if err != nil {
			http.Error(w, "内部エラー", http.StatusInternalServerError)
			return
		}

		proxy := &httputil.ReverseProxy{
			Rewrite: func(req *httputil.ProxyRequest) {
				req.Out.URL = targetURL
				req.Out.Host = targetURL.Host
				// 認証ヘッダー等をクライアントに漏らさないためにリセット
				req.Out.Header.Del("Authorization")
			},
		}
		proxy.ServeHTTP(w, r)
	})
}
