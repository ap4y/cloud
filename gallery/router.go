package gallery

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
)

type GalleryAPI struct {
	http.Handler
	source Source
	cache  Cache
}

func NewGalleryAPI(source Source, cache Cache) *GalleryAPI {
	mux := chi.NewRouter()
	api := &GalleryAPI{mux, source, cache}

	mux.Route("/", func(r chi.Router) {
		r.Get("/", api.listAlbums)
		r.Get("/{galleryName}", api.listAlbumImages)
		r.Get("/images/*", api.getImage)
		r.Get("/thumbnails/*", api.getImageThumbnail)
		r.Get("/exif/*", api.getImageEXIF)
	})

	return api
}

func (api *GalleryAPI) listAlbums(w http.ResponseWriter, req *http.Request) {
	albums, err := api.source.Albums()
	if err != nil {
		http.Error(w, fmt.Sprint("failed to fetch albums:", err), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(albums); err != nil {
		http.Error(w, fmt.Sprint("failed to encode albums:", err), http.StatusBadRequest)
	}
}

func (api *GalleryAPI) listAlbumImages(w http.ResponseWriter, req *http.Request) {
	galleryName := chi.URLParam(req, "galleryName")
	images, err := api.source.Images(galleryName)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to fetch images:", err), http.StatusBadRequest)
		return
	}

	if err := json.NewEncoder(w).Encode(images); err != nil {
		http.Error(w, fmt.Sprint("failed to encode albums:", err), http.StatusBadRequest)
	}
}

func (api *GalleryAPI) getImage(w http.ResponseWriter, req *http.Request) {
	imgPath := chi.URLParam(req, "*")
	file, err := api.source.Image(imgPath)
	if err != nil {
		log.Print("failed to fetch image:", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	fi, err := file.Stat()
	if err != nil {
		log.Print("failed to read image file stats:", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	http.ServeContent(w, req, imgPath, fi.ModTime(), file)
}

func (api *GalleryAPI) getImageThumbnail(w http.ResponseWriter, req *http.Request) {
	imgPath := chi.URLParam(req, "*")
	if thumb, modTime := api.cache.Thumbnail(imgPath); thumb != nil {
		http.ServeContent(w, req, imgPath, modTime, thumb)
		return
	}

	file, err := api.source.Image(imgPath)
	if err != nil {
		log.Print("failed to fetch image:", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	thumbData, err := ThumbnailForFile(file, 200)
	if err != nil {
		log.Print("failed to generate thumbnail:", err)
		http.Error(w, "", http.StatusNotFound)
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

func (api *GalleryAPI) getImageEXIF(w http.ResponseWriter, req *http.Request) {
	imgPath := chi.URLParam(req, "*")
	file, err := api.source.Image(imgPath)
	if err != nil {
		log.Print("failed to fetch image:", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	exif, err := EXIFForFile(file)
	if err != nil {
		log.Print("failed to generate thumbnail:", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	if err := json.NewEncoder(w).Encode(exif); err != nil {
		http.Error(w, fmt.Sprint("failed to encode exif data:", err), http.StatusBadRequest)
	}
}
