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
	fmt.Printf("Conn Type : %T\n", conn)
	fmt.Printf("Conn RemoteAddr : %s\n", conn.RemoteAddr())
	fmt.Printf("Conn LocalAddr : %s\n", conn.LocalAddr())

	buffer := make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("n : %d\n", n)
}