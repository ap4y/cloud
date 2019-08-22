package api

import (
	"context"
	"encoding/json"
	"io"
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

func TestShareAuthenticator(t *testing.T) {
	dir, err := ioutil.TempDir("", "shares")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store := NewDiskShareStore(dir)
	share := &Share{Slug: "bar", Type: ModuleGallery, Items: []string{"foo/test.jpg"}}
	require.NoError(t, store.Save(share))

	var ctxShare *Share
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctxShare, _ = r.Context().Value(ShareCtxKey).(*Share)
		io.WriteString(w, "Hello World!")
	}
	subMux := chi.NewRouter()
	subMux.Get("/{path}", handler)
	subMux.Get("/{path}/file/{file}", handler)

	mux := chi.NewRouter()
	mux.Route("/{slug}", func(r chi.Router) {
		r.Get("/root", handler)

		r.Group(func(r chi.Router) {
			r.Use(ShareAuthenticator(store))
			r.Mount("/folder", subMux)
		})
	})

	tcs := []struct {
		name   string
		path   string
		status int
		share  *Share
		body   string
	}{
		{"root", "/bar/root", http.StatusOK, nil, "Hello World!"},
		{"without path param", "/bar/folder", http.StatusNotFound, nil, "Not Found\n"},
		{"with path param", "/bar/folder/foo", http.StatusOK, share, "Hello World!"},
		{"unknown share", "/baz/folder/foo", http.StatusNotFound, nil, "Not Found\n"},
		{"unknown item", "/bar/folder/baz", http.StatusNotFound, share, "Not Found\n"},
		{"with file param", "/bar/folder/foo/file/test.jpg", http.StatusOK, share, "Hello World!"},
		{"unknown file", "/bar/folder/foo/file/baz.jpg", http.StatusNotFound, share, "Not Found\n"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			ctxShare = nil
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "http://cloud.api"+tc.path, nil)

			mux.ServeHTTP(w, req)
			resp := w.Result()
			require.Equal(t, tc.status, resp.StatusCode)

			assert.Equal(t, tc.share, ctxShare)
			body, _ := ioutil.ReadAll(resp.Body)
			assert.Equal(t, tc.body, string(body))
		})
	}
}

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
		ctx := context.WithValue(req.Context(), ShareCtxKey, share)
		handler.ServeHTTP(w, req.WithContext(ctx))

		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)

		resShare := &Share{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(resShare))
		require.NotNil(t, resShare.Slug)
	})
}
