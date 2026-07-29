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

	fmt.Printf("Listening on Addr : %s\n", listener.Addr())

	conn, err := listener.Accept()
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	fmt.Printf("Conn RemoteAddr : %s\n", conn.RemoteAddr())
	fmt.Printf("Conn LocalAddr : %s\n", conn.LocalAddr())
}