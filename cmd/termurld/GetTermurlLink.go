package main

import "net/http"

func (app *Application) GetTermurlLink(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	originalLink, err := app.db.GetOriginalLinkFromSlug(r.Context(), slug)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Slug Not Found")
		return
	}

	http.Redirect(w, r, originalLink, http.StatusFound)
}
