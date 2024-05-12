package main

import (
	"github.com/tuku13/image-gallery/api"
)

func main() {
	server := api.NewServer()
	server.Start()
}
