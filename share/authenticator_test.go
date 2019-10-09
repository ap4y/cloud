package share

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/ap4y/cloud/common"
)

func TestShareAuthenticator(t *testing.T) {
	dir, err := ioutil.TempDir("", "shares")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store, err := NewDiskStore(dir)
	require.NoError(t, err)
	share := &Share{Slug: "bar", Type: common.ModuleGallery, Name: "foo", Items: []string{"test.jpg"}}
	require.NoError(t, store.Save(share))

	var ctxShare *Share
	handler := func(w http.ResponseWriter, r *http.Request) {
		ctxShare, _ = r.Context().Value(common.ShareCtxKey).(*Share)
		io.WriteString(w, "Hello World!") // nolint: errcheck
	}

	mux := chi.NewRouter()
	mux.Route("/{slug}", func(r chi.Router) {
		r.Get("/root", handler)

		r.Group(func(r chi.Router) {
			r.Use(Authenticator(store))
			r.Get("/folder", handler)
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
		{"with path param", "/bar/folder", http.StatusOK, share, "Hello World!"},
		{"unknown share", "/baz/folder", http.StatusNotFound, nil, "Not Found\n"},
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
