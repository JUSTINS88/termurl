package main

import (
	"net/http"

	"github.com/JUSTINS88/zaps/internal/web"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", http.FileServerFS(web.Assets))
	mux.HandleFunc("GET /zaps-mark.svg", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, web.Assets, "zaps-mark.svg")
	})
	mux.HandleFunc("GET /stats", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFileFS(w, r, web.Assets, "stats.html")
	})
	mux.HandleFunc("POST /api/links", app.CreateNewLink)
	mux.HandleFunc("GET /api/links/{slug}/stats", app.GetLinksStats)
	mux.HandleFunc("GET /{slug}", app.GetTermurlLink)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Host == "s.justinsep.me" {
			http.Redirect(w, r, "https://zaps.my.id"+r.URL.RequestURI(), http.StatusMovedPermanently)
			return
		}
		mux.ServeHTTP(w, r)
	})
}
