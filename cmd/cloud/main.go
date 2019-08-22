package main

import (
	"flag"
	"log"

	"github.com/ap4y/cloud/internal/cli"
)

var (
	configPath = flag.String("config", "cloud.json", "path to a config file")
	addr       = flag.String("addr", ":8080", "address to server on")
	devURL     = flag.String("devURL", "", "url for a dev react web server")
)

func main() {
	flag.Parse()

	if err := cli.Run(*configPath, *devURL, *addr); err != nil {
		log.Fatal(err)
	}
}
