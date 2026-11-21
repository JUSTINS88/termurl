package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *Application) GetTermurlLink(w http.ResponseWriter, r *http.Request) {
	slug := r.PathValue("slug")

	// Check in cache
	if originalLink, ok := app.cache.Get(slug); ok {
		go app.trackClicks(slug)
		RespondWithPreview(w, slug, originalLink)
		return
	}

	// Not found in cache, search in DB
	originalLink, err := app.db.GetOriginalLinkFromSlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			RespondWithError(w, http.StatusNotFound, "Slug not found")
			return
		}
		app.logger.Error("Internal server error", "slug", slug, "error", err.Error())
		RespondWithError(w, http.StatusInternalServerError, "Server error")
		return
	}

	go app.trackClicks(slug)
	// Set cache after found in DB
	app.cache.Set(slug, originalLink)
	RespondWithPreview(w, slug, originalLink)
}

func (app *Application) increaseClick(ctx context.Context, slug string) error {
	rows, err := app.db.UpdateClicks(ctx, slug)
	if err != nil {
		return err
	}
	if rows == 0 {
		return fmt.Errorf("Slug not found")
	}
	return nil
}

func (app *Application) trackClicks(slug string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := app.increaseClick(ctx, slug); err != nil {
		app.logger.Error("failed to increase click", "slug", slug, "error", err)
	}
}
