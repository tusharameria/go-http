package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type PayloadType string

const (
	JOIN PayloadType = "JOIN"
	NAME PayloadType = "NAME"
	MSG  PayloadType = "MSG"
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
		conn.Write([]byte("Please join the the chatroom before sending messgaes\n"))
		conn.Write([]byte("JOIN <Your Name>\n"))
	}
}

func handleConnection(conn net.Conn, id uuid.UUID, connPool *ConnectionPool, msgInfoCh chan<- MessageInfo) {
	defer closeConnection(conn, id, connPool)

	buff := make([]byte, 8)
	var tmp bytes.Buffer
	for {
		// if err := conn.SetReadDeadline(time.Now().Add(10 * time.Second)); err != nil {
		// 	fmt.Printf("Error : %s", err)
		// 	return
		// }
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

		tmp.Write(buff[:n])

		for {
			data := tmp.Bytes()
			idx := strings.IndexByte(string(data), '\n')
			if idx == -1 {
				break
			}
			line := string(data[:idx])
			tmp.Next(idx + 1)

			if valid := validateProtocol(line, conn, id, connPool, msgInfoCh); !valid {
				fmt.Printf("Error validating protocol: %s\n", err)
				continue
			}
		}
	}
}

func validateProtocol(line string, conn net.Conn, id uuid.UUID, connPool *ConnectionPool, msgInfoCh chan<- MessageInfo) bool {
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, "\n")
	line = strings.TrimSuffix(line, "\r")
	parts := strings.SplitN(line, " ", 2)

	if line == "" || len(parts) == 0 {
		return false
	}
	if len(parts) != 2 {
		conn.Write([]byte(fmt.Sprintln("Invalid call!!!")))
		conn.Write([]byte(fmt.Sprintln("Here's the list of valid calls :")))
		conn.Write([]byte(fmt.Sprintln("JOIN <Your Name>")))
		conn.Write([]byte(fmt.Sprintln("NAME <New Name>")))
		conn.Write([]byte(fmt.Sprintln("MSG <Your Message>")))
		return false
	}
	callType := parts[0]
	callPayload := parts[1]

	fmt.Printf("Received from Client %s\n", id.String())
	fmt.Printf("Message : %s\n", callPayload)

	connPool.mu.Lock()
	switch PayloadType(callType) {
	case JOIN:
		if connPool.conns[id].name == nil {
			if strings.TrimSpace(callPayload) == "" {
				conn.Write([]byte(fmt.Sprintln("Please enter a valid name!!!")))
				return false
			}
			name := strings.TrimSpace(callPayload)
			connPool.conns[id].name = &name
		}
	case NAME:
		if connPool.conns[id].name == nil {
			conn.Write([]byte(fmt.Sprintln("Join the chatroom first!!!")))
			conn.Write([]byte(fmt.Sprintln("JOIN <Your Name>")))
		} else {
			if strings.TrimSpace(callPayload) == "" {
				conn.Write([]byte(fmt.Sprintln("Please enter a valid name!!!")))
				return false
			}
			name := strings.TrimSpace(callPayload)
			connPool.conns[id].name = &name
		}
	case MSG:
		if connPool.conns[id].name == nil {
			conn.Write([]byte(fmt.Sprintln("Join the chatroom first!!!")))
			conn.Write([]byte(fmt.Sprintln("JOIN <Your Name>")))
		} else {
			msgInfoCh <- MessageInfo{
				id:  id,
				msg: fmt.Sprintln(callPayload),
			}
		}
	default:
		conn.Write([]byte(fmt.Sprintln("Invalid call!!!")))
		conn.Write([]byte(fmt.Sprintln("Here's the list of valid calls :")))
		conn.Write([]byte(fmt.Sprintln("JOIN <Your Name>")))
		conn.Write([]byte(fmt.Sprintln("NAME <New Name>")))
		conn.Write([]byte(fmt.Sprintln("MSG <Your Message>")))
	}
	connPool.mu.Unlock()
	return true
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
