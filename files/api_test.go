package files

import (
	"bytes"
	"encoding/base64"
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

	t.Run("listDir", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/./files", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		items := make([]*Item, 0)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&items))
		require.Len(t, items, 3)
		assert.Equal(t, "foo", items[0].Name)
	})

	t.Run("getContents", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/./files/foo", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/base64", resp.Header.Get("Content-Type"))

		dec := base64.NewDecoder(base64.StdEncoding, resp.Body)
		data, err := ioutil.ReadAll(dec)
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
		req := httptest.NewRequest("POST", "http://cloud.api/./files", &buf)
		req.Header.Set("Content-Type", formWriter.FormDataContentType())
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		items, err := src.List(".")
		require.NoError(t, err)
		require.Len(t, items, 4)

		w = httptest.NewRecorder()
		req = httptest.NewRequest("DELETE", "http://cloud.api/./files/bar", nil)
		api.ServeHTTP(w, req)

		resp = w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		items, err = src.List(".")
		require.NoError(t, err)
		require.Len(t, items, 3)
	})
}
