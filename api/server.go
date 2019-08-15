package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// Module defines supported module
type Module string

const (
	// ModuleGallery represents gallery module
	ModuleGallery = "gallery"
)

// NewServer returns a new root handler for the app.
func NewServer(modules map[Module]http.Handler, cs CredentialsStorage, ss ShareStore) (http.Handler, error) {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	apiMux := chi.NewRouter()
	mux.Mount("/api", apiMux)

	apiMux.Mount("/user", AuthHandler(cs))
	apiMux.Get("/share/{slug}", getShareHandler(ss))
	apiMux.Group(func(r chi.Router) {
		r.Use(Authenticator(cs))

		r.Get("/modules", func(w http.ResponseWriter, res *http.Request) {
			moduleIds := []Module{}
			for module := range modules {
				moduleIds = append(moduleIds, module)
			}

			if err := json.NewEncoder(w).Encode(map[string][]Module{"modules": moduleIds}); err != nil {
				http.Error(w, fmt.Sprintf("Failed to encode json: %s", err), http.StatusBadRequest)
			}
		})

		r.Post("/share", createShareHandler(ss))

		for module, handler := range modules {
			r.Mount("/"+string(module), handler)
		}
	})

	return mux, nil
}
