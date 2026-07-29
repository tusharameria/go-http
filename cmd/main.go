package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8082")
	if err != nil {
		log.Panic(err)
	}
	defer listener.Close()

	fmt.Printf("Addr : %s\n", listener.Addr())
	fmt.Println(listener)
}