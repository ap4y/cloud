package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"

	"github.com/ap4y/cloud/api"
	"github.com/ap4y/cloud/gallery"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
)

type galleryConfig struct {
	Path  string `json:"path"`
	Cache string `json:"cache"`
}

type shareConfig struct {
	Path string `json:"path"`
}

type config struct {
	JWTSecret string            `json:"jwt_secret"`
	Modules   []api.Module      `json:"modules"`
	Users     map[string]string `json:"users"`
	Share     *shareConfig      `json:"share"`
	Gallery   *galleryConfig    `json:"gallery"`
}

var (
	configPath = flag.String("config", "cloud.json", "path to a config file")
	addr       = flag.String("addr", ":8080", "address to server on")
	devURL     = flag.String("devURL", "", "url for a dev react web server")
)

func main() {
	flag.Parse()

	cfg := new(config)
	f, err := os.Open(*configPath)
	if err != nil {
		log.Fatal("failed to open config file:", err)
	}

	if err := json.NewDecoder(f).Decode(cfg); err != nil {
		log.Fatal("failed to decode config file:", err)
	}

	srv, err := setupServer(cfg)
	if err != nil {
		log.Fatal("failed to initialise server:", err)
	}

	if *devURL != "" {
		setupDevProxy(*devURL, srv)
	}

	log.Println("Serving on", *addr)
	if err := http.ListenAndServe(*addr, srv); err != nil {
		log.Fatal("failed to start server:", err)
	}
}

func setupServer(cfg *config) (http.Handler, error) {
	modules := map[api.Module]http.Handler{}
	for _, module := range cfg.Modules {
		if module == "gallery" {
			gallery, err := galleryModule(cfg.Gallery)
			if err != nil {
				log.Fatal("failed to initialise gallery:", err)
				continue
			}

			modules[module] = gallery
		}
	}

	cs := api.NewMemoryCredentialsStorage(cfg.Users, jwt.SigningMethodHS256, []byte(cfg.JWTSecret))
	ss := api.NewDiskShareStore(cfg.Share.Path)

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

func setupDevProxy(devURL string, handler http.Handler) {
	rpURL, err := url.Parse(devURL)
	if err != nil {
		log.Fatal("invalid dev url:", err)
	}

	if mux, ok := handler.(*chi.Mux); ok {
		mux.Get("/", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
		mux.Get("/images/*", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
		mux.Get("/static/*", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
		mux.Get("/sockjs-node/*", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
	}
}

func galleryModule(cfg *galleryConfig) (http.Handler, error) {
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
