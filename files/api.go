package files

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/ap4y/cloud/internal/httputil"
	"github.com/go-chi/chi"
)

type apiItem struct {
	*Item
	Children []apiItem `json:"children"`
	URL      string    `json:"url"`
}

type filesAPI struct {
	http.Handler
	source Source
}

// NewFilesAPI returns a new http.Handler instance that implements files related endpoints.
func NewFilesAPI(source Source) http.Handler {
	mux := chi.NewRouter()
	api := &filesAPI{mux, source}

	mux.Route("/", func(r chi.Router) {
		r.Get("/", api.listTree)
		r.Post("/mkdir/{path}*", api.createFolder)
		r.Post("/upload/{path}*", api.uploadFile)
		r.Get("/file/{path}*", api.getFile)
		r.Delete("/file/{path}*", api.removeFile)
	})

	return api
}

func (api *filesAPI) listTree(w http.ResponseWriter, req *http.Request) {
	tree, err := api.source.Tree()
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to traverse path:", err), http.StatusBadRequest)
		return
	}

	httputil.Respond(w, toAPIItem(tree))
}

func (api *filesAPI) createFolder(w http.ResponseWriter, req *http.Request) {
	item, err := api.source.Mkdir(chi.URLParam(req, "path"))
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to created folder:", err), http.StatusBadRequest)
		return
	}

	httputil.Respond(w, toAPIItem(item))
}

func (api *filesAPI) uploadFile(w http.ResponseWriter, req *http.Request) {
	file, header, err := req.FormFile("file")
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to parse upload:", err), http.StatusBadRequest)
		return
	}

	fullPath := filepath.Join("/", chi.URLParam(req, "path"), header.Filename)
	path, err := url.QueryUnescape(fullPath)
	if err != nil {
		httputil.Error(w, fmt.Sprint("invalid path:", err), http.StatusBadRequest)
		return
	}

	defer file.Close()
	item, err := api.source.Save(file, path)
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to save upload:", err), http.StatusBadRequest)
		return
	}

	httputil.Respond(w, toAPIItem(item))
}

func (api *filesAPI) getFile(w http.ResponseWriter, req *http.Request) {
	fullPath := filepath.Join("/", chi.URLParam(req, "path"))
	path, err := url.QueryUnescape(fullPath)
	if err != nil {
		httputil.Error(w, fmt.Sprint("invalid path:", err), http.StatusBadRequest)
		return
	}

	file, err := api.source.File(path)
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to traverse path:", err), http.StatusBadRequest)
		return
	}

	fi, err := file.Stat()
	if err != nil {
		http.Error(w, fmt.Sprint("failed to read file stats:", err), http.StatusNotFound)
		return
	}

	http.ServeContent(w, req, fullPath, fi.ModTime(), file)
}

func (api *filesAPI) removeFile(w http.ResponseWriter, req *http.Request) {
	fullPath := filepath.Join("/", chi.URLParam(req, "path"))
	path, err := url.QueryUnescape(fullPath)
	if err != nil {
		httputil.Error(w, fmt.Sprint("invalid path:", err), http.StatusBadRequest)
		return
	}

	item, err := api.source.Remove(path)
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to remove file:", err), http.StatusBadRequest)
		return
	}

	httputil.Respond(w, toAPIItem(item))
}

func itemPath(item *Item) string {
	if item.Type == ItemTypeFile {
		path := strings.ReplaceAll(item.Path, item.Name, "")
		if path == "/" {
			path = "//"
		}

		return fmt.Sprintf("%sfile/%s", path, item.Name)
	}

	return item.Path
}

func toAPIItem(item *Item) apiItem {
	aItem := apiItem{Item: item, URL: itemPath(item)}
	aItem.Children = apiTree(item.Children)
	return aItem
}

func apiTree(tree []*Item) []apiItem {
	result := make([]apiItem, len(tree))
	for idx, item := range tree {
		result[idx] = toAPIItem(item)
	}

	return result
}
