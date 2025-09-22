package main

import (
	"fmt"
	"net"
	"os"
	"sync"

	"github.com/google/uuid"
)

type ConnectionPool struct {
	conns map[uuid.UUID]net.Conn
	mu    sync.Mutex
}

type MessageInfo struct {
	id  uuid.UUID
	msg string
}

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:8082")
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Listener Speifications : %+v\n", listener)
	fmt.Printf("Listener Address : %+v\n", listener.Addr())

	msgInfoCh := make(chan MessageInfo)

	connPool := &ConnectionPool{
		conns: make(map[uuid.UUID]net.Conn),
		mu:    sync.Mutex{},
	}

	go func() {
		for msgInfo := range msgInfoCh {
			broadcast(msgInfo.msg, msgInfo.id, connPool)
		}
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("%s\n", err)
			continue
		}
		newID := addToConnections(conn, connPool)
		fmt.Println("New Client Connected...")
		fmt.Printf("Local address : %+v\n", conn.LocalAddr())
		fmt.Printf("Active Connections : %d\n", len(connPool.conns))

		go handleConnection(conn, newID, connPool, msgInfoCh)
		conn.Write([]byte("Hello from minimal TCP server...\n"))
	}
}

func handleConnection(conn net.Conn, id uuid.UUID, connPool *ConnectionPool, msgInfoCh chan<- MessageInfo) {
	defer closeConnection(conn, id, connPool)

	buff := make([]byte, 1024)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("%s\n", err)
			return
		}
		fmt.Printf("Received from Client %s\n", id.String())
		fmt.Printf("Message : %s", buff[:n])
		msgInfoCh <- MessageInfo{
			id:  id,
			msg: string(buff[:n]),
		}
	}
}

func closeConnection(conn net.Conn, id uuid.UUID, connPool *ConnectionPool) {
	conn.Close()
	removeFromConnections(id, connPool)
}

func addToConnections(conn net.Conn, connPool *ConnectionPool) uuid.UUID {
	connPool.mu.Lock()
	defer connPool.mu.Unlock()

	id := uuid.New()
	connPool.conns[id] = conn
	return id
}

func removeFromConnections(id uuid.UUID, connPool *ConnectionPool) {
	connPool.mu.Lock()
	defer connPool.mu.Unlock()

	delete(connPool.conns, id)
	fmt.Printf("Closing Connections : %s\n", id.String())
	fmt.Printf("Active Connections : %d\n", len(connPool.conns))
}

func broadcast(msg string, senderID uuid.UUID, connPool *ConnectionPool) {
	connPool.mu.Lock()
	defer connPool.mu.Unlock()

	for id, conn := range connPool.conns {
		if id != senderID {
			conn.Write([]byte(fmt.Sprintf("Client %s : %s", senderID.String(), msg)))
		}
	}
}
