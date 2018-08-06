package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

func main() {
	mux := chi.NewRouter()
	mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
	})

	if err := http.ListenAndServe(":3000", mux); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
