package gallery

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ap4y/cloud/common"
	"github.com/ap4y/cloud/internal/httputil"
	"github.com/ap4y/cloud/share"
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
		r.Get("/", share.BlockHandler(api.listAlbums))
		r.Route("/{gallery}", func(r chi.Router) {
			r.Get("/images", share.VerifyHandler(common.ModuleGallery, "gallery", "", api.listAlbumImages))
			r.Get("/image/{file}", share.VerifyHandler(common.ModuleGallery, "gallery", "file", api.getImage))
			r.Get("/thumbnail/{file}", share.VerifyHandler(common.ModuleGallery, "gallery", "file", api.getImageThumbnail))
			r.Get("/exif/{file}", share.VerifyHandler(common.ModuleGallery, "gallery", "file", api.getImageEXIF))
		})
	})

	return api
}

func (api *galleryAPI) listAlbums(w http.ResponseWriter, req *http.Request) {
	albums, err := api.source.Albums()
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to fetch albums:", err), http.StatusBadRequest)
		return
	}

	httputil.Respond(w, albums)
}

func (api *galleryAPI) listAlbumImages(w http.ResponseWriter, req *http.Request) {
	galleryName := chi.URLParam(req, "gallery")
	images, err := api.source.Images(galleryName)
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to fetch images:", err), http.StatusBadRequest)
		return
	}

	if share, ok := req.Context().Value(common.ShareCtxKey).(*share.Share); ok {
		shareImages := make([]Image, 0, len(images))
		for _, image := range images {
			if share.Includes(galleryName, image.Path) {
				shareImages = append(shareImages, image)
			}
		}
		images = shareImages
	}

	httputil.Respond(w, images)
}

func (api *galleryAPI) getImage(w http.ResponseWriter, req *http.Request) {
	file, err := api.source.Image(chi.URLParam(req, "gallery"), chi.URLParam(req, "file"))
	if err != nil {
		http.Error(w, fmt.Sprint("failed to fetch image:", err), http.StatusNotFound)
		return
	}

	fi, err := file.Stat()
	if err != nil {
		http.Error(w, fmt.Sprint("failed to read image file stats:", err), http.StatusNotFound)
		return
	}

	http.ServeContent(w, req, fi.Name(), fi.ModTime(), file)
}

func (api *galleryAPI) getImageThumbnail(w http.ResponseWriter, req *http.Request) {
	galleryName := chi.URLParam(req, "gallery")
	fileName := chi.URLParam(req, "file")

	if thumb, modTime := api.cache.Thumbnail(galleryName, fileName); thumb != nil {
		http.ServeContent(w, req, fileName, modTime, thumb)
		return
	}

	file, err := api.source.Image(galleryName, fileName)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to fetch image:", err), http.StatusNotFound)
		return
	}

	thumbData, err := Thumbnail(file, 200)
	if err != nil {
		http.Error(w, fmt.Sprint("failed to generate thumbnail:", err), http.StatusNotFound)
		return
	}

	thumb, err := api.cache.StoreThumbnail(galleryName, fileName, thumbData)
	if err != nil {
		log.Print("failed to cache thumbnail:", err)
		http.Error(w, "", http.StatusNotFound)
		return
	}

	http.ServeContent(w, req, fileName, time.Now(), thumb)
}

func (api *galleryAPI) getImageEXIF(w http.ResponseWriter, req *http.Request) {
	file, err := api.source.Image(chi.URLParam(req, "gallery"), chi.URLParam(req, "file"))
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to fetch image:", err), http.StatusNotFound)
		return
	}

	exif, err := EXIF(file)
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to parse exif:", err), http.StatusNotFound)
		return
	}

	httputil.Respond(w, exif)
}
