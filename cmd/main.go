package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Listener Speifications : %+v\n", listener)
	fmt.Printf("Listener Address : %+v\n", listener.Addr())

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s\n", err)
		}

		fmt.Printf("New client connected, local address : %+v\n", connection.LocalAddr())
		fmt.Printf("New client connected, remote address : %+v\n", connection.RemoteAddr())

		connection.Write([]byte("Hello from minimal TCP server\n"))

		connection.Close()
	}
}
