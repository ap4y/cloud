package files

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/ap4y/cloud/common"
	"gitlab.com/ap4y/cloud/contextkey"
	"gitlab.com/ap4y/cloud/share"
)

func TestFilesAPI(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	src, err := NewDiskSource(filepath.Join(pwd, "fixtures"))
	require.NoError(t, err)

	api := NewFilesAPI(src)

	share := &share.Share{Type: common.ModuleFiles, Name: "/test1", Items: []string{"/test1/inner"}}

	t.Run("listTree", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		tree := &apiItem{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(tree))
		assert.Equal(t, "/", tree.Name)
		assert.Equal(t, "/", tree.URL)
	})

	t.Run("listTree/with share", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/", nil)
		ctx := context.WithValue(req.Context(), contextkey.ShareCtxKey, share)
		api.ServeHTTP(w, req.WithContext(ctx))

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		tree := &apiItem{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(tree))
		assert.Equal(t, "test1", tree.Name)
		assert.Equal(t, "/test1", tree.Path)
		require.Len(t, tree.Children, 1)
		item := tree.Children[0]
		assert.Equal(t, "inner", item.Name)
		assert.Equal(t, "/test1/inner", item.Path)
		require.Len(t, item.Children, 1)
	})

	t.Run("getFile", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/file/foo", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "foo\n", string(data))
	})

	t.Run("getFile/with share", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/file/test1/inner/foo", nil)
		ctx := context.WithValue(req.Context(), contextkey.ShareCtxKey, share)
		api.ServeHTTP(w, req.WithContext(ctx))

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "foo\n", string(data))
	})

	t.Run("getFile/with unmatched share", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/file/foo", nil)
		ctx := context.WithValue(req.Context(), contextkey.ShareCtxKey, share)
		api.ServeHTTP(w, req.WithContext(ctx))

		resp := w.Result()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("uploadFile/removeFile", func(t *testing.T) {
		var buf bytes.Buffer
		formWriter := multipart.NewWriter(&buf)
		fw, err := formWriter.CreateFormFile("file", "bar")
		require.NoError(t, err)
		_, err = fw.Write([]byte("bar"))
		require.NoError(t, err)
		formWriter.Close()

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "http://cloud.api/upload/test1/inner", &buf)
		req.Header.Set("Content-Type", formWriter.FormDataContentType())
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		tree, err := src.Tree()
		require.NoError(t, err)
		require.Len(t, tree.Children, 3)

		w = httptest.NewRecorder()
		req = httptest.NewRequest("DELETE", "http://cloud.api/file/test1/inner/bar", nil)
		api.ServeHTTP(w, req)

		resp = w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		tree, err = src.Tree()
		require.NoError(t, err)
		require.Len(t, tree.Children, 3)
	})

	t.Run("uploadFile/with share", func(t *testing.T) {
		var buf bytes.Buffer
		formWriter := multipart.NewWriter(&buf)
		fw, err := formWriter.CreateFormFile("file", "bar")
		require.NoError(t, err)
		_, err = fw.Write([]byte("bar"))
		require.NoError(t, err)
		formWriter.Close()

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "http://cloud.api/upload/test1/inner", &buf)
		req.Header.Set("Content-Type", formWriter.FormDataContentType())
		ctx := context.WithValue(req.Context(), contextkey.ShareCtxKey, share)
		api.ServeHTTP(w, req.WithContext(ctx))

		resp := w.Result()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("removeFile/with share", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "http://cloud.api/file/test1/inner/foo", nil)
		ctx := context.WithValue(req.Context(), contextkey.ShareCtxKey, share)
		api.ServeHTTP(w, req.WithContext(ctx))

		resp := w.Result()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("createFolder/removeFolder", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "http://cloud.api/mkdir/test1/subfolder", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		item := &apiItem{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(item))
		assert.Equal(t, "subfolder", item.Name)
		assert.Equal(t, "/test1/subfolder", item.URL)

		w = httptest.NewRecorder()
		req = httptest.NewRequest("POST", "http://cloud.api/rmdir/test1/subfolder", nil)
		api.ServeHTTP(w, req)

		resp = w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		item = &apiItem{}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(item))
		assert.Equal(t, "subfolder", item.Name)
		assert.Equal(t, "/test1/subfolder", item.URL)
	})

	t.Run("createFolder/with share", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "http://cloud.api/mkdir/test1/subfolder", nil)
		ctx := context.WithValue(req.Context(), contextkey.ShareCtxKey, share)
		api.ServeHTTP(w, req.WithContext(ctx))

		resp := w.Result()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("removeFolder/with share", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "http://cloud.api/rmdir/test1/subfolder", nil)
		ctx := context.WithValue(req.Context(), contextkey.ShareCtxKey, share)
		api.ServeHTTP(w, req.WithContext(ctx))

		resp := w.Result()
		require.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
