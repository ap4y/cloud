package api

import (
	"net/http"

	"github.com/ap4y/cloud/internal/httputil"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type contextKey struct {
	name string
}

// UsernameCtxKey defines username request context key.
var UsernameCtxKey = &contextKey{"Username"}

// ShareCtxKey defines share request context key.
var ShareCtxKey = &contextKey{"Share"}

// Module defines supported module
type Module string

const (
	// ModuleGallery represents gallery module
	ModuleGallery = "gallery"

	// ModuleFiles represents files module
	ModuleFiles = "files"
)

// NewServer returns a new root handler for the app.
func NewServer(modules map[Module]http.Handler, cs CredentialsStorage, ss ShareStore) (http.Handler, error) {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	sh := &shareHandler{ss}
	mux.Route("/api", func(apiMux chi.Router) {
		apiMux.Mount("/user", AuthHandler(cs))

		apiMux.Group(func(r chi.Router) {
			r.Use(Authenticator(cs))

			moduleIds := make([]Module, len(modules))
			idx := 0
			for module := range modules {
				moduleIds[idx] = module
				idx++
			}

			r.Get("/modules", func(w http.ResponseWriter, res *http.Request) {
				httputil.Respond(w, map[string][]Module{"modules": moduleIds})
			})

			r.Get("/shares", sh.listShares)
			r.Post("/shares", sh.createShare)
			r.Delete("/shares/{slug}", sh.removeShare)

			for module, handler := range modules {
				r.Mount("/"+string(module), handler)
			}
		})

		apiMux.Route("/share/{slug}", func(r chi.Router) {
			r.Get("/", sh.getShare)

			r.Group(func(r chi.Router) {
				r.Use(ShareAuthenticator(ss))

				for module, handler := range modules {
					r.Mount("/"+string(module), handler)
				}
			})
		})
	})

	return mux, nil
}
