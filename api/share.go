package api

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ap4y/cloud/internal/httputil"
	"github.com/go-chi/chi"
)

// Share stores share data.
type Share struct {
	Slug      string    `json:"slug"`
	Type      Module    `json:"type"`
	Name      string    `json:"name"`
	Items     []string  `json:"items"`
	ExpiresAt time.Time `json:"expires_at"`
}

// ShareStore manages share metadata.
type ShareStore interface {
	// Save persists share metadata.
	Save(share *Share) error
	// Get return share metadata.
	Get(slug string) (*Share, error)
	// Remove removes share metadata.
	Remove(slug string) error
	// Expire removes all expired shares.
	Expire() error
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

func (store *diskShareStore) Remove(slug string) error {
	path := filepath.Join(store.dir, slug)
	return os.Remove(path)
}

func (store *diskShareStore) Expire() error {
	path := filepath.Join(store.dir, "*")
	matches, err := filepath.Glob(path)
	if err != nil {
		return fmt.Errorf("file: %s", err)
	}

	for _, match := range matches {
		_, slug := filepath.Split(match)
		share, err := store.Get(slug)
		if err != nil {
			continue
		}

		if share.ExpiresAt.After(time.Now()) {
			continue
		}

		if err := store.Remove(slug); err != nil {
			return err
		}
	}

	return nil
}

type wrapResponseWriter struct {
	http.ResponseWriter
	buf        bytes.Buffer
	statusCode int
}

func (w *wrapResponseWriter) Write(buf []byte) (int, error) {
	return w.buf.Write(buf)
}

func (w *wrapResponseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
}

// ShareAuthenticator returns new share authentication middleware.
func ShareAuthenticator(store ShareStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			slug := chi.URLParam(req, "slug")
			if slug == "" {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			share, err := store.Get(slug)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			ctx := context.WithValue(req.Context(), ShareCtxKey, share)
			wrapper := &wrapResponseWriter{ResponseWriter: w}
			next.ServeHTTP(wrapper, req.WithContext(ctx))

			path := chi.URLParam(req, "path")
			if path == "" {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}

			file := chi.URLParam(req, "file")
			fullPath := filepath.Join(path, file)
			for _, item := range share.Items {
				if strings.HasPrefix(item, fullPath) {
					w.Write(wrapper.buf.Bytes())
					if wrapper.statusCode != 0 {
						w.WriteHeader(wrapper.statusCode)
					}
					return
				}
			}

			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		})
	}
}

func getShareHandler(store ShareStore) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		slug := chi.URLParam(req, "slug")
		if slug == "" {
			httputil.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		share, err := store.Get(slug)
		if err != nil {
			httputil.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		httputil.Respond(w, share)
	}
}

func createShareHandler(store ShareStore) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var share *Share
		if err := json.NewDecoder(req.Body).Decode(&share); err != nil {
			httputil.Error(w, fmt.Sprintf("Failed to decode json: %s", err), http.StatusBadRequest)
			return
		}

		valid := false
		if share.Type == ModuleGallery {
			valid = share.Name != ""
		}

		if !valid {
			httputil.Error(w, "Invalid share", http.StatusUnprocessableEntity)
			return
		}

		slug := make([]byte, 10)
		if _, err := rand.Read(slug); err != nil {
			httputil.Error(w, fmt.Sprintf("Failed to generate slug: %s", err), http.StatusBadRequest)
			return
		}

		share.Slug = base64.StdEncoding.EncodeToString(slug)
		if err := store.Save(share); err != nil {
			httputil.Error(w, fmt.Sprintf("Failed to save: %s", err), http.StatusBadRequest)
			return
		}

		httputil.Respond(w, share)
	}
}
