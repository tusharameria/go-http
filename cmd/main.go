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

	bufferLen := 8
	buffer := make([]byte, bufferLen)
	msg := ""
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			break
		}
		msg = string(buffer[:n])
		fmt.Printf("Msg Received : %s\n", msg)

		conn.Write([]byte(fmt.Sprintf("Response from server : Hey you sent us \"%s\"\n", msg)))
	}
}
