package files

import (
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"

	"gitlab.com/ap4y/cloud/common"
	"gitlab.com/ap4y/cloud/internal/httputil"
	"gitlab.com/ap4y/cloud/share"
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
		r.Get("/", verifyHandler("", api.listTree))
		r.Post("/mkdir/{path}*", share.BlockHandler(api.createFolder))
		r.Post("/rmdir/{path}*", share.BlockHandler(api.removeFolder))
		r.Post("/upload/{path}*", share.BlockHandler(api.uploadFile))
		r.Get("/file/{path}*", verifyHandler("path", api.getFile))
		r.Delete("/file/{path}*", share.BlockHandler(api.removeFile))
	})

	return api
}

func (api *filesAPI) listTree(w http.ResponseWriter, req *http.Request) {
	tree, err := api.source.Tree()
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to traverse path:", err), http.StatusBadRequest)
		return
	}

	if share, ok := req.Context().Value(common.ShareCtxKey).(*share.Share); ok {
		tree = locateTreeNode(tree, share.Name)
		if tree == nil {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		filtered := make([]*Item, 0, len(tree.Children))
		for _, child := range tree.Children {
			if share.Includes(tree.Path, child.Path) {
				filtered = append(filtered, child)
			}
		}
		tree.Children = filtered
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

func (api *filesAPI) removeFolder(w http.ResponseWriter, req *http.Request) {
	item, err := api.source.Rmdir(chi.URLParam(req, "path"))
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
	filePath := chi.URLParam(req, "path")
	file, err := api.source.File(filePath)
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to traverse path:", err), http.StatusBadRequest)
		return
	}

	fi, err := file.Stat()
	if err != nil {
		http.Error(w, fmt.Sprint("failed to read file stats:", err), http.StatusNotFound)
		return
	}

	http.ServeContent(w, req, fi.Name(), fi.ModTime(), file)
}

func (api *filesAPI) removeFile(w http.ResponseWriter, req *http.Request) {
	filePath := chi.URLParam(req, "path")
	item, err := api.source.Remove(filePath)
	if err != nil {
		httputil.Error(w, fmt.Sprint("failed to remove file:", err), http.StatusBadRequest)
		return
	}

	httputil.Respond(w, toAPIItem(item))
}

func itemPath(item *Item) string {
	if item.Type == ItemTypeFile {
		return fmt.Sprintf("/file%s", item.Path)
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

func locateTreeNode(tree *Item, path string) *Item {
	if path == "/" {
		return tree
	}

	node := tree
	components := strings.Split(path, "/")
	for _, component := range components {
		if component == "" {
			continue
		}

		var item *Item
		for _, child := range node.Children {
			if child.Name == component {
				item = child
				break
			}
		}

		if item == nil {
			return nil
		}

		node = item
	}

	return node
}

func verifyHandler(itemParam string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		share, ok := req.Context().Value(common.ShareCtxKey).(*share.Share)
		if !ok {
			next.ServeHTTP(w, req)
			return
		}

		if share.Type != common.ModuleFiles {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if itemParam == "" {
			next.ServeHTTP(w, req)
			return
		}

		filePath := "/" + chi.URLParam(req, itemParam)
		for _, item := range share.Items {
			if strings.HasPrefix(filePath, item) {
				next.ServeHTTP(w, req)
				return
			}
		}

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
}
