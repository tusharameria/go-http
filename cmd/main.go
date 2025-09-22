package main

import (
	"fmt"
	"net"
	"os"
	"sync/atomic"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Listener Speifications : %+v\n", listener)
	fmt.Printf("Listener Address : %+v\n", listener.Addr())

	var connCount int64 = 0

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		atomic.AddInt64(&connCount, 1)
		fmt.Println("New Client Connected...")
		fmt.Printf("Local address : %+v\n", conn.LocalAddr())
		fmt.Printf("Active Connections : %d\n", connCount)

		go handleConnection(conn, &connCount)
		conn.Write([]byte("Hello from minimal TCP server...\n"))
	}
}

func handleConnection(conn net.Conn, connCount *int64) {
	defer closeConnection(conn, connCount)

	buff := make([]byte, 1024)
	num := *connCount
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Received from Client %d\n", num)
		fmt.Printf("Message : %s", buff[:n])

		conn.Write([]byte(fmt.Sprintf("You sent : %s", buff[:n])))
	}
}

func closeConnection(conn net.Conn, connCount *int64) {
	newVal := atomic.AddInt64(connCount, -1)
	fmt.Printf("Active Connections after closing connection : %d\n", newVal)
	conn.Close()
}
