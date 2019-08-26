package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/ap4y/cloud/api"
	"github.com/ap4y/cloud/gallery"
	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
)

// Run is an entry point for a CLI.
func Run(configPath, devURL, addr string) error {
	cfg := new(Config)
	f, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to open config file: %s", err)
	}

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		return fmt.Errorf("failed to decode config file: %s", err)
	}

	srv, err := setupServer(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialise server: %s", err)
	}

	if devURL != "" {
		if err := setupDevProxy(devURL, srv); err != nil {
			return err
		}
	}

	log.Println("Serving on", addr)

	if err := http.ListenAndServe(addr, srv); err != nil {
		return fmt.Errorf("failed to start server: %s", err)
	}

	return nil
}

func setupServer(cfg *Config) (http.Handler, error) {
	modules := map[api.Module]http.Handler{}
	for _, module := range cfg.Modules {
		if module == "gallery" {
			gallery, err := galleryModule(cfg.Gallery)
			if err != nil {
				return nil, fmt.Errorf("failed to initialise gallery: %s", err)
			}

			modules[module] = gallery
		}
	}

	cs := api.NewMemoryCredentialsStorage(cfg.Users, jwt.SigningMethodHS256, []byte(cfg.JWTSecret))
	ss, err := api.NewDiskShareStore(cfg.Share.Path)
	if err != nil {
		return nil, fmt.Errorf("Failed to create share store: %s", err)
	}

	expireTicker := time.NewTicker(time.Hour)
	go func() {
		for range expireTicker.C {
			if err := ss.Expire(); err != nil {
				log.Println("failed to expire shares:", err)
			}
		}
	}()

	return api.NewServer(modules, cs, ss)
}

func setupDevProxy(devURL string, handler http.Handler) error {
	rpURL, err := url.Parse(devURL)
	if err != nil {
		return fmt.Errorf("invalid dev url: %s", err)
	}

	if mux, ok := handler.(*chi.Mux); ok {
		mux.Get("/*", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
	}

	return nil
}

func galleryModule(cfg *GalleryConfig) (http.Handler, error) {
	source, err := gallery.NewDiskSource(cfg.Path, []string{".jpg", ".jpeg", ".png"})
	if err != nil {
		return nil, err
	}

	cache, err := gallery.NewDiskCache(cfg.Cache)
	if err != nil {
		return nil, err
	}

	return gallery.NewGalleryAPI(source, cache), nil
}
