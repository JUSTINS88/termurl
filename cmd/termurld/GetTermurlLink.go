package main

import "net/http"

func (app *Application) GetTermurlLink(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	// Check in cache
	if originalLink, ok := app.cache.Get(slug); ok {
		http.Redirect(w, r, originalLink, http.StatusFound)
		return
	}

	// Not found in cache, search in DB
	originalLink, err := app.db.GetOriginalLinkFromSlug(r.Context(), slug)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Slug Not Found")
		return
	}

	// Set cache after found in DB
	app.cache.Set(slug, originalLink)
	http.Redirect(w, r, originalLink, http.StatusFound)
}
