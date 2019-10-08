package httputil

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Error encodes json error payload into ResponseWriter.
func Error(w http.ResponseWriter, error string, code int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": error}) // nolint: errcheck
}

// Respond json encodes entity into ResponseWriter or error if encoding failed.
func Respond(w http.ResponseWriter, entity interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(entity); err != nil {
		Error(w, fmt.Sprint("failed to encode entity:", err), http.StatusBadRequest)
	}
}
