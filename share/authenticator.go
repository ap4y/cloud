package share

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ap4y/cloud/contextkey"
)

// Authenticator returns new share authentication middleware.
func Authenticator(store Store) func(next http.Handler) http.Handler {
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

			ctx := context.WithValue(req.Context(), contextkey.ShareCtxKey, share)
			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}
