package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ConnectionInfo struct {
	name *string
	conn net.Conn
}

type ConnectionPool struct {
	conns map[uuid.UUID]*ConnectionInfo
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
		conns: make(map[uuid.UUID]*ConnectionInfo),
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
		conn.Write([]byte("Please enter you name\n"))
	}
}

func handleConnection(conn net.Conn, id uuid.UUID, connPool *ConnectionPool, msgInfoCh chan<- MessageInfo) {
	defer closeConnection(conn, id, connPool)

	buff := make([]byte, 1024)
	for {
		if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
			fmt.Printf("Error : %s", err)
			return
		}
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %s disconnected\n", id.String())
			} else {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					fmt.Printf("Client %s timed out\n", id.String())
					conn.Write([]byte("You have been timed out due to inactivity.\n"))
					return
				}
				fmt.Printf("Read error: %s\n", err)
			}
			return
		}

		msg := string(buff[:n])

		fmt.Printf("Received from Client %s\n", id.String())
		fmt.Printf("Message : %s", msg)

		if msg == "\n" {
			conn.Write([]byte(fmt.Sprintln("Please enter a valid name!!!")))
			continue
		}

		connPool.mu.Lock()
		if connPool.conns[id].name != nil {
			msgInfoCh <- MessageInfo{
				id:  id,
				msg: msg,
			}
		} else {
			name := strings.TrimSpace(msg)
			connPool.conns[id].name = &name
		}
		connPool.mu.Unlock()
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
	connPool.conns[id] = &ConnectionInfo{
		name: nil,
		conn: conn,
	}
	return id
}

func removeFromConnections(id uuid.UUID, connPool *ConnectionPool) {
	connPool.mu.Lock()
	defer connPool.mu.Unlock()

	delete(connPool.conns, id)
	fmt.Printf("Closing Connection : %s\n", id.String())
	fmt.Printf("Active Connections : %d\n", len(connPool.conns))
}

func broadcast(msg string, senderID uuid.UUID, connPool *ConnectionPool) {
	connPool.mu.Lock()
	defer connPool.mu.Unlock()

	for id, connInfo := range connPool.conns {
		if id != senderID && connInfo.name != nil {
			connInfo.conn.Write([]byte(fmt.Sprintf("%s : %s", *connPool.conns[senderID].name, msg)))
		}
	}
}
