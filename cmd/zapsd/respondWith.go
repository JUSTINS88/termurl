package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"github.com/JUSTINS88/zaps/internal/web"
)

var previewTmpl = template.Must(template.ParseFS(web.Assets, "preview.html"))

func RespondWithPreview(w http.ResponseWriter, slug, originalLink string) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if err := previewTmpl.Execute(w, struct {
		Slug         string
		OriginalLink string
	}{
		Slug:         slug,
		OriginalLink: originalLink,
	}); err != nil {
		log.Printf("Error rendering preview page: %v", err)
	}
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, struct {
		Error string `json:"error"`
	}{Error: message})
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")

	msg, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Error marshaling json: %v", err)
		return
	}

	w.WriteHeader(code)
	w.Write(msg)
}
