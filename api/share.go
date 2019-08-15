package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
)

// TODO: expirations
// TODO: module integration (unrestrict content)

// Share stores share data.
type Share struct {
	Slug  string   `json:"slug"`
	Type  Module   `json:"type"`
	Items []string `json:"items"`
}

// ShareStore manages share metadata.
type ShareStore interface {
	// Save persists share metadata.
	Save(share *Share) error
	// Get return share metadata.
	Get(slug string) (*Share, error)
}

type diskShareStore struct {
	dir string
}

// NewDiskShareStore returns a new on-disk implementation of the ShareStore.
func NewDiskShareStore(dir string) ShareStore {
	return &diskShareStore{dir}
}

func (store *diskShareStore) Save(share *Share) error {
	path := filepath.Join(store.dir, share.Slug)
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return fmt.Errorf("file: %s", err)
	}

	if err := json.NewEncoder(file).Encode(share); err != nil {
		return fmt.Errorf("json: %s", err)
	}

	return nil
}

func (store *diskShareStore) Get(slug string) (*Share, error) {
	path := filepath.Join(store.dir, slug)
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if err != nil {
		return nil, fmt.Errorf("file: %s", err)
	}

	share := &Share{}
	if err := json.NewDecoder(file).Decode(share); err != nil {
		return nil, fmt.Errorf("json: %s", err)
	}

	return share, nil
}

func getShareHandler(store ShareStore) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		slug := chi.URLParam(req, "slug")
		share, err := store.Get(slug)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(share); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode json: %s", err), http.StatusBadRequest)
		}
	}
}

func createShareHandler(store ShareStore) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var share *Share
		w.Header().Set("Content-Type", "application/json")

		if err := json.NewDecoder(req.Body).Decode(&share); err != nil {
			http.Error(w, fmt.Sprintf("Failed to decode json: %s", err), http.StatusBadRequest)
			return
		}

		slug := make([]byte, 10)
		if _, err := rand.Read(slug); err != nil {
			http.Error(w, fmt.Sprintf("Failed to generate slug: %s", err), http.StatusBadRequest)
			return
		}

		share.Slug = base64.StdEncoding.EncodeToString(slug)
		if err := store.Save(share); err != nil {
			http.Error(w, fmt.Sprintf("Failed to save: %s", err), http.StatusBadRequest)
			return
		}

		if err := json.NewEncoder(w).Encode(share); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode json: %s", err), http.StatusBadRequest)
		}
	}
}
