package main

import (
	"log"

	"github.com/tusharameria/go-http/internal/server"
)

func main() {
	srv, err := server.New(":8082")
	if err != nil {
		log.Fatal(err)
	}

	log.Fatal(srv.Serve())
}
