package main

import (
	"fmt"
	"log"
	"net"

	"github.com/tusharameria/go-http/internal/request"
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

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)

		fmt.Printf("Headers:\n")
		r.Headers.ForEach(func(k, v string) {
			fmt.Printf("- %s: %s\n", k, v)
		})

		fmt.Printf("Body:\n")
		fmt.Printf("%s\n", r.Body)
	}
}
