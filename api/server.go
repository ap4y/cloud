package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"

	"github.com/ap4y/cloud/gallery"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	http.Handler
}

func NewServer(modules []string, galleryPath string, devURL string) (*Server, error) {
	mux := chi.NewRouter()

	mux.Use(middleware.Logger)

	// API routes
	apiMux := chi.NewRouter()
	mux.Mount("/api", apiMux)
	apiMux.Get("/modules", func(w http.ResponseWriter, red *http.Request) {
		if err := json.NewEncoder(w).Encode(map[string][]string{"modules": modules}); err != nil {
			http.Error(w, fmt.Sprintf("Failed to encode json: %s", err), http.StatusBadRequest)
		}
	})

	for _, module := range modules {
		switch module {
		case "gallery":
			source, err := gallery.NewDiskSource(galleryPath, []string{".jpg", ".jpeg", ".png"})
			if err != nil {
				return nil, fmt.Errorf("failed to initialise gallery: %s", err)
			}

			cache, err := gallery.NewDiskCache(filepath.Join(os.TempDir(), "cloud"))
			if err != nil {
				return nil, fmt.Errorf("failed to initialise gallery: %s", err)
			}

			apiMux.Mount("/gallery", gallery.NewGalleryAPI(source, cache))
		default:
			return nil, fmt.Errorf("unknown module: %s", module)
		}
	}

	// Static files
	if devURL != "" {
		rpURL, err := url.Parse(devURL)
		if err != nil {
			return nil, fmt.Errorf("invalid dev url: %s", err)
		}
		mux.Get("/", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
		mux.Get("/images/*", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
		mux.Get("/static/*", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
		// mux.Get("/sockjs-node/*", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
		// mux.Post("/sockjs-node/*", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
		// mux.Get("/__webpack_dev_server__/*", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
	}

	return &Server{mux}, nil
}
