package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShareAuthenticator(t *testing.T) {
	dir, err := ioutil.TempDir("", "shares")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store, err := NewDiskShareStore(dir)
	require.NoError(t, err)
	share := &Share{Slug: "bar", Type: ModuleGallery, Name: "foo", Items: []string{"test.jpg"}}
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

func TestShareStore(t *testing.T) {
	dir, err := ioutil.TempDir("", "shares")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store, err := NewDiskShareStore(dir)
	require.NoError(t, err)

	share := &Share{Slug: "foo", Type: ModuleGallery, Items: []string{"foo", "bar"}, ExpiresAt: time.Time{}}
	t.Run("Save", func(t *testing.T) {
		require.NoError(t, store.Save(share))
	})

	t.Run("All", func(t *testing.T) {
		res, err := store.All()
		require.NoError(t, err)
		require.Len(t, res, 1)
		assert.Equal(t, share, res[0])
	})

	t.Run("Get", func(t *testing.T) {
		res, err := store.Get("foo")
		require.NoError(t, err)
		assert.Equal(t, share, res)
	})

	t.Run("Remove", func(t *testing.T) {
		require.NoError(t, store.Remove("foo"))

		res, err := store.Get("foo")
		require.Error(t, err)
		assert.Nil(t, res)
	})

	t.Run("Expire", func(t *testing.T) {
		require.NoError(t, store.Save(share))
		require.NoError(t, store.Expire())

		res, err := store.Get("foo")
		require.Error(t, err)
		assert.Nil(t, res)
	})
}

func TestShareHandler(t *testing.T) {
	dir, err := ioutil.TempDir("", "shares")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store, err := NewDiskShareStore(dir)
	require.NoError(t, err)

	sh := &shareHandler{store}
	handler := chi.NewRouter()
	handler.Get("/{slug}", sh.getShare)
	handler.Delete("/{slug}", sh.removeShare)
	handler.Post("/", sh.createShare)
	handler.Get("/", sh.listShares)

	share := &Share{Slug: "foo", Type: ModuleGallery, Items: []string{"foo", "bar"}, ExpiresAt: time.Time{}}
	require.NoError(t, store.Save(share))

	t.Run("List", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/", nil)
		handler.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)

		shares := make([]*Share, 0)
		require.NoError(t, json.NewDecoder(res.Body).Decode(&shares))
		require.Len(t, shares, 1)
	})

	t.Run("Fetch", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/foo", nil)
		handler.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)

		resShare := &Share{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(resShare))
		require.NotNil(t, resShare.Slug)
	})

	t.Run("Remove", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("DELETE", "http://cloud.api/foo", nil)
		handler.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)

		shares, err := store.All()
		require.NoError(t, err)
		assert.Len(t, shares, 0)
	})

	t.Run("Create", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := "{\"type\":\"gallery\",\"name\":\"test\",\"items\":[\"foo\",\"bar\"],\"expires_at\":\"1970-01-01T00:00:00Z\"}"
		req := httptest.NewRequest("POST", "http://cloud.api/", strings.NewReader(body))
		handler.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)

		share := &Share{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(share))

		require.NotNil(t, share.Slug)
		assert.Equal(t, ModuleGallery, string(share.Type))
		assert.Equal(t, int64(0), share.ExpiresAt.Unix())
	})

	t.Run("Create - invalid", func(t *testing.T) {
		w := httptest.NewRecorder()
		body := "{\"type\":\"gallery\",\"items\":[\"foo\",\"bar\"],\"expires_at\":\"1970-01-01T00:00:00Z\"}"
		req := httptest.NewRequest("POST", "http://cloud.api/", strings.NewReader(body))
		handler.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusUnprocessableEntity, res.StatusCode)
	})
}
