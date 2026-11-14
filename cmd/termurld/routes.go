package main

import (
	"net/http"

	"github.com/JUSTINS88/termurl/internal/web"
)

func (app *Application) Routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /{$}", http.FileServerFS(web.Assets))
	mux.HandleFunc("POST /api/links", app.CreateNewLink)
	mux.HandleFunc("GET /{slug}", app.GetTermurlLink)

	return mux
}
