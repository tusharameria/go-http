package server

import (
	"fmt"
	"net"

	"github.com/tusharameria/go-http/internal/parser"
)

func (s *Server) handleConnection(conn net.Conn) error {
	bufferLen := 8
	buffer := make([]byte, bufferLen)
	p := parser.Parser{}
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return err
		}
		res := p.Feed(buffer[:n])
		for _, line := range res {
			fmt.Printf("Received : %s\n", line)
			conn.Write([]byte(fmt.Sprintf("You sent : %s\n", line)))
		}
	}
}
