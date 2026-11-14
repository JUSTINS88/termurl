package main

import (
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
)

func sendLinkToCreateLink(originalLink, slug string) error {
	machineKey, err := GetMachineKey()
	if err != nil {
		return fmt.Errorf("Machine Key cannot be generated: %v", err)
	}

	link, err := json.Marshal(struct {
		OriginalLink string `json:"original_link"`
		Slug         string `json:"slug"`
		MachineKey   string `json:"machine_key"`
	}{
		OriginalLink: originalLink,
		Slug:         slug,
		MachineKey:   machineKey,
	})

	res, err := http.Post("http://localhost:8080/api/links", "application/json", bytes.NewBuffer(link))
	if err != nil {
		return fmt.Errorf("Failed Connect to server")
	}

	defer res.Body.Close()

	// Check if status Created
	if res.StatusCode != http.StatusCreated {
		var errResp struct {
			Error string `json:"error"`
		}
		if err := json.NewDecoder(res.Body).Decode(&errResp); err == nil && errResp.Error != "" {
			return fmt.Errorf("server error (%d): %s", res.StatusCode, errResp.Error)
		}
		return fmt.Errorf("server returned status %d", res.StatusCode)
	}

	// Show in terminal
	var result struct {
		ShortURL string `json:"short_url"`
	}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return fmt.Errorf("Failed To Decode JSON")
	}

	fmt.Printf("This is your shorten link: %s\n", result.ShortURL)
	return nil
}

func GetMachineKey() (string, error) {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".termurl.json")

	// Get Machine Key
	if data, err := os.ReadFile(path); err == nil {
		var cfg struct {
			MachineKey string `json:"machine_key"`
		}
		if json.Unmarshal(data, &cfg) == nil && cfg.MachineKey != "" {
			return cfg.MachineKey, nil
		}
	}

	// No Machine Key
	newMachineKey := rand.Text()
	data, _ := json.Marshal(map[string]string{
		"machine_key": newMachineKey,
	})
	_ = os.WriteFile(path, data, 0600)
	return newMachineKey, nil
}
