package gallery

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ap4y/cloud/common"
	"github.com/ap4y/cloud/share"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGalleryAPI(t *testing.T) {
	pwd, err := os.Getwd()
	require.NoError(t, err)

	src, err := NewDiskSource(filepath.Join(pwd, "fixtures"), []string{".jpg"})
	require.NoError(t, err)

	dir, err := ioutil.TempDir("", "cache")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	cache, err := NewDiskCache(dir)
	require.NoError(t, err)

	api := NewGalleryAPI(src, cache)

	t.Run("listAlbums", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		albums := make([]*Album, 0)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&albums))
		require.Len(t, albums, 2)
		assert.Equal(t, "album1", albums[0].Name)
	})

	t.Run("listAlbumImages", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/album1/images", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		images := make([]*Image, 0)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&images))
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		require.Len(t, images, 1)
		assert.Equal(t, "test", images[0].Name)
	})

	t.Run("listAlbumImages/with_share", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/album1/images", nil)
		share := &share.Share{Type: common.ModuleGallery, Name: "album1", Items: []string{"test.jpg"}}
		ctx := context.WithValue(req.Context(), common.ShareCtxKey, share)
		api.ServeHTTP(w, req.WithContext(ctx))

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		images := make([]*Image, 0)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&images))
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		require.Len(t, images, 1)
		assert.Equal(t, "test", images[0].Name)
	})

	t.Run("listAlbumImages/with_share/no_match", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/album1/images", nil)
		share := &share.Share{Type: common.ModuleGallery, Name: "album1", Items: []string{"foo.jpg"}}
		ctx := context.WithValue(req.Context(), common.ShareCtxKey, share)
		api.ServeHTTP(w, req.WithContext(ctx))

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)

		images := make([]*Image, 0)
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&images))
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))

		require.Len(t, images, 0)
	})

	t.Run("getImage", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/album1/image/test.jpg", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"))
	})

	t.Run("getImageThumbnail", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/album1/thumbnail/test.jpg", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "image/jpeg", resp.Header.Get("Content-Type"))
	})

	t.Run("getImageEXIF", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/album1/exif/test.jpg", nil)
		api.ServeHTTP(w, req)

		resp := w.Result()
		require.Equal(t, http.StatusOK, resp.StatusCode)
		assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	})
}
