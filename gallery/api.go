package gallery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type galleryAPI struct {
	http.Handler
	source Source
	cache  Cache
}

// NewGalleryAPI returns a new http.Handler instance that implements gallery related endpoints.
func NewGalleryAPI(source Source, cache Cache) http.Handler {
	mux := chi.NewRouter()
	api := &galleryAPI{mux, source, cache}

	mux.Route("/", func(r chi.Router) {
		r.Get("/", api.listAlbums)
		r.Get("/{galleryName}/images", api.listAlbumImages)
		r.Get("/images/*", api.getImage)
		r.Get("/thumbnails/*", api.getImageThumbnail)
		r.Get("/exif/*", api.getImageEXIF)
	})

	return api
}

func (api *galleryAPI) listAlbums(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	albums, err := api.source.Albums()
	if err != nil {
		http.Error(w, fmt.Sprint("failed to fetch albums:", err), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(albums); err != nil {
		http.Error(w, fmt.Sprint("failed to encode albums:", err), http.StatusBadRequest)
	}
}

func (api *galleryAPI) listAlbumImages(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	galleryName := chi.URLParam(req, "galleryName")
	images, err := api.source.Images(galleryName)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to fetch images:", err), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(images); err != nil {
		http.Error(w, fmt.Sprint("failed to encode images:", err), http.StatusBadRequest)
	}
}

func (api *galleryAPI) getImage(w http.ResponseWriter, req *http.Request) {
	imgPath := chi.URLParam(req, "*")
	file, err := api.source.Image(imgPath)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to fetch image:", err), http.StatusNotFound)
		return
	}

	fi, err := file.Stat()
	if err != nil {
		http.Error(w, fmt.Sprint("failed to read image file stats:", err), http.StatusNotFound)
		return
	}

	http.ServeContent(w, req, imgPath, fi.ModTime(), file)
}

func (api *galleryAPI) getImageThumbnail(w http.ResponseWriter, req *http.Request) {
	imgPath := chi.URLParam(req, "*")
	if thumb, modTime := api.cache.Thumbnail(imgPath); thumb != nil {
		http.ServeContent(w, req, imgPath, modTime, thumb)
		return
	}

	file, err := api.source.Image(imgPath)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to fetch image:", err), http.StatusNotFound)
		return
	}

	thumbData, err := Thumbnail(file, 200)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to generate thumbnail:", err), http.StatusNotFound)
		return
	}

	thumb, err := api.cache.StoreThumbnail(imgPath, thumbData)
	if err != nil {
		log.Print("failed to cache thumbnail:", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	http.ServeContent(w, req, imgPath, time.Now(), thumb)
}

func (api *galleryAPI) getImageEXIF(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	imgPath := chi.URLParam(req, "*")
	file, err := api.source.Image(imgPath)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to fetch image:", err), http.StatusNotFound)
		return
	}

	exif, err := EXIF(file)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to parse exif:", err), http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(exif); err != nil {
		http.Error(w, fmt.Sprint("failed to encode exif data:", err), http.StatusBadRequest)
	}
}
