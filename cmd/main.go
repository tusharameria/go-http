package main

import (
	"bytes"
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

	buffer := make([]byte, 1024)
	for {
		_, err := conn.Read(buffer)
		if err != nil {
			break
		}
		newLineIdx := bytes.LastIndex(buffer, []byte("\n"))
		if newLineIdx == -1 {
			continue
		}
		msg := string(buffer[:newLineIdx])
		fmt.Printf("Msg Received : %s\n", msg)

		conn.Write([]byte(fmt.Sprintf("Response from server : Hey you sent us \"%s\"\n", msg)))
	}
}
