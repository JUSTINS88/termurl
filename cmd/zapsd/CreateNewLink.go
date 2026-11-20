package main

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"regexp"

	"github.com/JUSTINS88/zaps/internal/store"
	"github.com/lib/pq"
)

const slugCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var validSlug = regexp.MustCompile(`^[a-zA-Z0-9_-]{1,50}$`)

type CreateLinkParams struct {
	OriginalLink string `json:"original_link"`
	Slug         string `json:"slug"`
	MachineKey   string `json:"machine_key"`
}

func (app *Application) CreateNewLink(w http.ResponseWriter, r *http.Request) {
	var createLink CreateLinkParams
	if err := json.NewDecoder(r.Body).Decode(&createLink); err != nil {
		RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Machine Key is valid or not
	if createLink.MachineKey == "" {
		RespondWithError(w, http.StatusBadRequest, "Machine key is missing")
		return
	}

	machineKey, err := app.db.CheckMachineKey(r.Context(), createLink.MachineKey)
	if err != nil {
		RespondWithError(w, http.StatusInternalServerError, "Failed to verify machine key")
		return
	}

	// Checking if original link exist
	if createLink.OriginalLink == "" {
		RespondWithError(w, http.StatusBadRequest, "Original link is required")
		return
	}

	// Checking if Slug exist
	if createLink.Slug == "" {
		randomSlug, err := generateRandomSlug()
		if err != nil {
			RespondWithError(w, http.StatusInternalServerError, "Failed to generate Slug")
			return
		}
		createLink.Slug = randomSlug
	}

	if createLink.Slug == "stats" || createLink.Slug == "api" {
		RespondWithError(w, http.StatusConflict, "You cannot use this slug")
		return
	}

	// Check is slug char is valid
	if !validSlug.MatchString(createLink.Slug) {
		RespondWithError(w, http.StatusBadRequest, "Slug can only contain letters, numbers, '-' and '_' and within 1-50 characters")
		return
	}

	// Insert into database
	createLinkRes, err := app.db.CreateNewLink(r.Context(), store.CreateNewLinkParams{
		OriginalLink: createLink.OriginalLink,
		Slug:         createLink.Slug,
		MachineKey:   machineKey,
	})
	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			RespondWithError(w, http.StatusConflict, "Slug already taken")
			return
		}
		app.logger.Error("failed to create new link", "error", err)
		RespondWithError(w, http.StatusInternalServerError, "Create Link failed.")
		return
	}

	// Shorten Link Success
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	shortURL := fmt.Sprintf("%s://%s/%s", scheme, r.Host, createLinkRes.Slug)

	RespondWithJSON(w, http.StatusCreated, struct {
		Short_url string `json:"short_url"`
	}{
		Short_url: shortURL,
	})
}

func generateRandomSlug() (string, error) {
	res := make([]byte, 6)
	for i := range res {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(slugCharset))))
		if err != nil {
			return "", err
		}
		res[i] = slugCharset[num.Int64()]
	}
	return string(res), nil
}
