package api

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/ap4y/cloud/internal/httputil"
	"github.com/ap4y/cloud/share"
)

type shareHandler struct {
	store share.Store
}

func (sh shareHandler) listShares(w http.ResponseWriter, req *http.Request) {
	shares, err := sh.store.All()
	if err != nil {
		httputil.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	httputil.Respond(w, shares)
}

func (sh shareHandler) getShare(w http.ResponseWriter, req *http.Request) {
	slug := chi.URLParam(req, "slug")
	if slug == "" {
		httputil.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	share, err := sh.store.Get(slug)
	if err != nil {
		httputil.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	httputil.Respond(w, share)
}

func (sh shareHandler) createShare(w http.ResponseWriter, req *http.Request) {
	var share *share.Share
	if err := json.NewDecoder(req.Body).Decode(&share); err != nil {
		httputil.Error(w, fmt.Sprintf("Failed to decode json: %s", err), http.StatusBadRequest)
		return
	}

	slug := make([]byte, 10)
	if _, err := rand.Read(slug); err != nil {
		httputil.Error(w, fmt.Sprintf("Failed to generate slug: %s", err), http.StatusBadRequest)
		return
	}

	share.Slug = base64.URLEncoding.EncodeToString(slug)
	if !share.IsValid() {
		httputil.Error(w, "Invalid share", http.StatusUnprocessableEntity)
		return
	}

	if err := sh.store.Save(share); err != nil {
		httputil.Error(w, fmt.Sprintf("Failed to save: %s", err), http.StatusBadRequest)
		return
	}

	httputil.Respond(w, share)
}

func (sh shareHandler) removeShare(w http.ResponseWriter, req *http.Request) {
	slug := chi.URLParam(req, "slug")
	if slug == "" {
		httputil.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	if err := sh.store.Remove(slug); err != nil {
		httputil.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	httputil.Respond(w, map[string]string{})
}
