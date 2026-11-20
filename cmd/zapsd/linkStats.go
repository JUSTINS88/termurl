package main

import (
	"database/sql"
	"errors"
	"net/http"
	"time"
)

func (app *Application) GetLinksStats(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")
	clicks, err := app.db.GetNumberOfClicks(r.Context(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, http.StatusNotFound, "Slug not found")
			return
		}
		app.logger.Error("Internal server error", "slug", slug, "error", err.Error())
		RespondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	RespondWithJSON(w, http.StatusOK, struct {
		OriginalLink string    `json:"original_link"`
		Slug         string    `json:"slug"`
		Clicks       int       `json:"clicks"`
		CreatedAt    time.Time `json:"created_at"`
	}{
		OriginalLink: clicks.OriginalLink,
		Slug:         clicks.Slug,
		Clicks:       int(clicks.Clicks),
		CreatedAt:    clicks.CreatedAt,
	})
}
