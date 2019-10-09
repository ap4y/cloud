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

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"gitlab.com/ap4y/cloud/api"
	"gitlab.com/ap4y/cloud/app"
	"gitlab.com/ap4y/cloud/common"
	"gitlab.com/ap4y/cloud/files"
	"gitlab.com/ap4y/cloud/gallery"
	"gitlab.com/ap4y/cloud/share"
)

// TODO: config validation

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

	if err := setupAssets(devURL, srv); err != nil {
		return err
	}

	log.Println("Serving on", addr)

	if err := http.ListenAndServe(addr, srv); err != nil {
		return fmt.Errorf("failed to start server: %s", err)
	}

	return nil
}

func setupServer(cfg *Config) (http.Handler, error) {
	modules := map[common.ModuleType]http.Handler{}
	for _, module := range cfg.Modules {
		var handler http.Handler
		var err error

		if module == common.ModuleGallery {
			handler, err = galleryModule(cfg.Gallery)
			if err != nil {
				return nil, fmt.Errorf("failed to initialise gallery: %s", err)
			}
		}

		if module == common.ModuleFiles {
			handler, err = filesModule(cfg.Files)
			if err != nil {
				return nil, fmt.Errorf("failed to initialise gallery: %s", err)
			}
		}

		modules[module] = handler
	}

	cs := api.NewMemoryCredentialsStorage(cfg.Users, jwt.SigningMethodHS256, []byte(cfg.JWTSecret))
	ss, err := share.NewDiskStore(cfg.Share.Path)
	if err != nil {
		return nil, fmt.Errorf("failed to create share store: %s", err)
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

func setupAssets(devURL string, handler http.Handler) error {
	mux, ok := handler.(*chi.Mux)
	if !ok {
		return fmt.Errorf("unsupported handler")
	}

	if devURL != "" {
		rpURL, err := url.Parse(devURL)
		if err != nil {
			return fmt.Errorf("invalid dev url: %s", err)
		}

		mux.Get("/*", httputil.NewSingleHostReverseProxy(rpURL).ServeHTTP)
	} else {
		fs := app.FS(false)
		mux.Handle("/", http.FileServer(fs))
		mux.Handle("/static/*", http.FileServer(fs))
		mux.NotFound(func(w http.ResponseWriter, req *http.Request) {
			file, _ := fs.Open("/index.html")
			stat, _ := file.Stat()
			http.ServeContent(w, req, "index.html", stat.ModTime(), file)
		})
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

func filesModule(cfg *FilesConfig) (http.Handler, error) {
	source, err := files.NewDiskSource(cfg.Path)
	if err != nil {
		return nil, err
	}

	return files.NewFilesAPI(source), nil
}
