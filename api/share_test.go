package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShareHandler(t *testing.T) {
	dir, err := ioutil.TempDir("", "shares")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store := NewDiskShareStore(dir)
	handler := chi.NewRouter()
	handler.Get("/{slug}", getShareHandler(store))
	handler.Post("/", createShareHandler(store))

	t.Run("Create", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := "{\"type\":\"gallery\",\"items\":[\"foo\",\"bar\"]}"
		req := httptest.NewRequest("POST", "http://cloud.api/", strings.NewReader(body))
		handler.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)

		share := &Share{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(share))

		require.NotNil(t, share.Slug)
		assert.Equal(t, ModuleGallery, string(share.Type))
	})

	t.Run("Fetch", func(t *testing.T) {
		share := &Share{Slug: "foo", Type: ModuleGallery, Items: []string{"foo", "bar"}}
		require.NoError(t, store.Save(share))

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/foo", nil)
		handler.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)

		resShare := &Share{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(resShare))
		require.NotNil(t, resShare.Slug)
	})
}
