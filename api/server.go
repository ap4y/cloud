package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

// NewServer returns a new root handler for the app.
func NewServer(modules map[string]http.Handler, cs CredentialsStorage) (http.Handler, error) {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	apiMux := chi.NewRouter()
	mux.Mount("/api", apiMux)

	apiMux.Mount("/user", AuthHandler(cs))
	apiMux.Group(func(r chi.Router) {
		r.Use(Authenticator(cs))

		r.Get("/modules", func(w http.ResponseWriter, res *http.Request) {
			moduleIds := []string{}
			for module := range modules {
				moduleIds = append(moduleIds, module)
			}

			if err := json.NewEncoder(w).Encode(map[string][]string{"modules": moduleIds}); err != nil {
				http.Error(w, fmt.Sprintf("Failed to encode json: %s", err), http.StatusBadRequest)
			}
		})

		for module, handler := range modules {
			r.Mount("/"+module, handler)
		}
	})

	return mux, nil
}
