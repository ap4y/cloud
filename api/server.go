package api

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"gitlab.com/ap4y/cloud/common"
	"gitlab.com/ap4y/cloud/internal/httputil"
	"gitlab.com/ap4y/cloud/share"
)

// NewServer returns a new root handler for the app.
func NewServer(modules map[common.ModuleType]http.Handler, cs CredentialsStorage, ss share.Store) (http.Handler, error) {
	mux := chi.NewRouter()
	mux.Use(middleware.Logger)

	sh := &shareHandler{ss}
	mux.Route("/api", func(apiMux chi.Router) {
		apiMux.Mount("/user", AuthHandler(cs))

		apiMux.Group(func(r chi.Router) {
			r.Use(Authenticator(cs))

			moduleIds := make([]common.ModuleType, len(modules))
			idx := 0
			for module := range modules {
				moduleIds[idx] = module
				idx++
			}

			r.Get("/modules", func(w http.ResponseWriter, res *http.Request) {
				httputil.Respond(w, map[string][]common.ModuleType{"modules": moduleIds})
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
				r.Use(share.Authenticator(ss))

				for module, handler := range modules {
					r.Mount("/"+string(module), handler)
				}
			})
		})
	})

	return mux, nil
}
