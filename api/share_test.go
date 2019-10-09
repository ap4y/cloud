package api

import (
	"encoding/json"
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

	"gitlab.com/ap4y/cloud/module"
	"gitlab.com/ap4y/cloud/niltime"
	"gitlab.com/ap4y/cloud/share"
)

func TestShareHandler(t *testing.T) {
	dir, err := ioutil.TempDir("", "shares")
	require.NoError(t, err)
	defer os.RemoveAll(dir)

	store, err := share.NewDiskStore(dir)
	require.NoError(t, err)

	sh := &shareHandler{store}
	handler := chi.NewRouter()
	handler.Get("/{slug}", sh.getShare)
	handler.Delete("/{slug}", sh.removeShare)
	handler.Post("/", sh.createShare)
	handler.Get("/", sh.listShares)

	s := &share.Share{
		Slug:      "foo",
		Type:      module.Gallery,
		Items:     []string{"foo", "bar"},
		ExpiresAt: niltime.Time{Time: time.Time{}},
	}
	require.NoError(t, store.Save(s))

	t.Run("List", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/", nil)
		handler.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)

		shares := make([]*share.Share, 0)
		require.NoError(t, json.NewDecoder(res.Body).Decode(&shares))
		require.Len(t, shares, 1)
	})

	t.Run("Fetch", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "http://cloud.api/foo", nil)
		handler.ServeHTTP(w, req)

		res := w.Result()
		require.Equal(t, http.StatusOK, res.StatusCode)

		resShare := &share.Share{}
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

		share := &share.Share{}
		require.NoError(t, json.NewDecoder(res.Body).Decode(share))

		require.NotNil(t, share.Slug)
		assert.Equal(t, module.Gallery, share.Type)
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
