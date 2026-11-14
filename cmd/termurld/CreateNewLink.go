package main

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/JUSTINS88/termurl/internal/store"
)

const slugCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

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
		RespondWithError(w, http.StatusInternalServerError, "Machine Key not found")
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

	// Insert into database
	createLinkRes, err := app.db.CreateNewLink(r.Context(), store.CreateNewLinkParams{
		OriginalLink: createLink.OriginalLink,
		Slug:         createLink.Slug,
		MachineKey:   machineKey,
	})
	if err != nil {
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
