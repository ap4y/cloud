package files

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"github.com/ap4y/cloud/internal/httputil"
	"github.com/go-chi/chi"
)

type filesAPI struct {
	http.Handler
	source Source
}

// NewFilesAPI returns a new http.Handler instance that implements files related endpoints.
func NewFilesAPI(source Source) http.Handler {
	mux := chi.NewRouter()
	api := &filesAPI{mux, source}

	mux.Route("/", func(r chi.Router) {
		r.Route("/{path}", func(r chi.Router) {
			r.Get("/files", api.listDir)
			r.Post("/files", api.uploadFile)
			r.Get("/files/{file}", api.getContents)
			r.Delete("/files/{file}", api.removeFile)
		})
	})

	return api
}

func (api *filesAPI) listDir(w http.ResponseWriter, req *http.Request) {
	path := chi.URLParam(req, "path")
	tree, err := api.source.List(path)
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to traverse path:", err), http.StatusBadRequest)
		return
	}

	httputil.Respond(w, tree)
}

func (api *filesAPI) uploadFile(w http.ResponseWriter, req *http.Request) {
	path := chi.URLParam(req, "path")
	file, header, err := req.FormFile("file")
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to parse upload:", err), http.StatusBadRequest)
		return
	}

	defer file.Close()
	if err := api.source.Save(file, path, header.Filename); err != nil {
		httputil.Error(w, fmt.Sprint("failed to save upload:", err), http.StatusBadRequest)
		return
	}

	httputil.Respond(w, map[string]string{})
}

func (api *filesAPI) getContents(w http.ResponseWriter, req *http.Request) {
	filePath := filepath.Join(chi.URLParam(req, "path"), chi.URLParam(req, "file"))
	contents, err := api.source.Content(filePath)
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to traverse path:", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/base64")

	enc := base64.NewEncoder(base64.StdEncoding, w)
	if _, err := io.Copy(enc, contents); err != nil {
		httputil.Error(w, fmt.Sprint("failed to encode file content:", err), http.StatusBadRequest)
	}

	contents.Close()
	enc.Close()
}

func (api *filesAPI) removeFile(w http.ResponseWriter, req *http.Request) {
	filePath := filepath.Join(chi.URLParam(req, "path"), chi.URLParam(req, "file"))
	if err := api.source.Remove(filePath); err != nil {
		httputil.Error(w, fmt.Sprint("failed to remove file:", err), http.StatusBadRequest)
		return
	}

	httputil.Respond(w, map[string]string{})
}
