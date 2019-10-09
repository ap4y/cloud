package share

import (
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/ap4y/cloud/common"
)

func TestBlockHandler(t *testing.T) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "<html><body>Hello World!</body></html>") // nolint: errcheck
	}

	t.Run("with share", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()

		share := &Share{Name: "album1", Items: []string{"test.jpg"}}
		ctx := context.WithValue(req.Context(), common.ShareCtxKey, share)
		BlockHandler(handler)(w, req.WithContext(ctx))

		resp := w.Result()
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("without share", func(t *testing.T) {
		req := httptest.NewRequest("GET", "http://example.com/foo", nil)
		w := httptest.NewRecorder()

		BlockHandler(handler)(w, req)

		resp := w.Result()
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	})
}

func TestVerifyHandler(t *testing.T) {
	share := &Share{Slug: "bar", Type: common.ModuleGallery, Name: "foo", Items: []string{"test.jpg"}}
	filesShare := &Share{Slug: "baz", Type: common.ModuleFiles, Name: "foo", Items: []string{"test.jpg"}}

	handler := func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World!") // nolint: errcheck
	}

	mux := chi.NewRouter()
	mux.Get("/{path}", VerifyHandler(common.ModuleGallery, "path", "", handler))
	mux.Get("/{path}/file/{file}", VerifyHandler(common.ModuleGallery, "path", "file", handler))

	tcs := []struct {
		name   string
		path   string
		status int
		share  *Share
		body   string
	}{
		{"without share", "/bar", http.StatusOK, nil, "Hello World!"},
		{"with path param", "/foo", http.StatusOK, share, "Hello World!"},
		{"unknown name", "/bar", http.StatusNotFound, share, "Not Found\n"},
		{"with file param", "/foo/file/test.jpg", http.StatusOK, share, "Hello World!"},
		{"unknown file", "/foo/file/baz.jpg", http.StatusNotFound, share, "Not Found\n"},
		{"invalid type", "/foo", http.StatusNotFound, filesShare, "Not Found\n"},
	}

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", "http://cloud.api"+tc.path, nil)
			w := httptest.NewRecorder()

			if tc.share != nil {
				ctx := context.WithValue(req.Context(), common.ShareCtxKey, tc.share)
				req = req.WithContext(ctx)
			}
			mux.ServeHTTP(w, req)

			resp := w.Result()
			require.Equal(t, tc.status, resp.StatusCode)
			body, err := ioutil.ReadAll(resp.Body)
			require.NoError(t, err)
			assert.Equal(t, tc.body, string(body))
		})
	}
}
