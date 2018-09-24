package main

import (
	"log"
	"net/http"

	"github.com/ap4y/cloud/api"
)

var (
	modules = []string{"gallery"}
	devURL  = "http://localhost:3000"
	addr    = ":8080"

	galleryPath = "/mnt/media/Photos/Export/2018/"
)

func main() {
	srv, err := api.NewServer(modules, galleryPath, devURL)
	if err != nil {
		log.Fatal("failed to initialise server:", err)
	}

	if err := http.ListenAndServe(addr, srv); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
