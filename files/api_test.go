package files

import (
	"bytes"
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
)

func TestGalleryAPI(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	src, err := NewDiskSource(filepath.Join(pwd, "fixtures"))
	require.NoError(t, err)

	api := NewFilesAPI(src)

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

	t.Run("getFile", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api//file/foo", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "text/plain; charset=utf-8", resp.Header.Get("Content-Type"))

		data, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		assert.Equal(t, "foo\n", string(data))
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
		req := httptest.NewRequest("POST", "http://cloud.api/./file", &buf)
		req.Header.Set("Content-Type", formWriter.FormDataContentType())
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		tree, err := src.Tree()
		require.NoError(t, err)
		require.Len(t, tree.Children, 4)

		w = httptest.NewRecorder()
		req = httptest.NewRequest("DELETE", "http://cloud.api/./file/bar", nil)
		api.ServeHTTP(w, req)

		resp = w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		tree, err = src.Tree()
		require.NoError(t, err)
		require.Len(t, tree.Children, 3)
	})
}
