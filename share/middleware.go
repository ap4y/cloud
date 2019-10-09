package share

import (
	"net/http"

	"github.com/go-chi/chi"

	"gitlab.com/ap4y/cloud/contextkey"
	"gitlab.com/ap4y/cloud/module"
)

// BlockHandler responds with NotFound for all requests that have share in context
func BlockHandler(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if _, ok := req.Context().Value(contextkey.ShareCtxKey).(*Share); ok {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		next.ServeHTTP(w, req)
	})
}

// VerifyHandler verifies share from the context against name and option item parameter.
func VerifyHandler(shareType module.Type, nameParam, itemParam string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		share, ok := req.Context().Value(contextkey.ShareCtxKey).(*Share)
		if !ok {
			next.ServeHTTP(w, req)
			return
		}

		if shareType != share.Type {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if chi.URLParam(req, nameParam) != share.Name {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}

		if itemParam == "" {
			next.ServeHTTP(w, req)
			return
		}

		file := chi.URLParam(req, itemParam)
		for _, item := range share.Items {
			if item == file {
				next.ServeHTTP(w, req)
				return
			}
		}

		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	})
}
