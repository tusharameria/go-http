package main

import (
	"fmt"
	"log"
	"net"

	"github.com/tusharameria/go-http/cmd/request"
)

func main() {
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "err", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("error in Accept", "err", err)
		}
		r, err := request.RequestFromReader(conn)
		if err != nil {
			log.Println("error in RequestFromReader", "err", err)
		}
		fmt.Println("Request Line :")
		fmt.Printf("- Method: %v\n", r.RequestLine.Method)
		fmt.Printf("- Target: %v\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %v\n", r.RequestLine.HttpVersion)
	}
}
