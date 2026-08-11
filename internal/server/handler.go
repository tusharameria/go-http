package server

import (
	"net"

	"github.com/tusharameria/go-http/internal/http"
)

func (s *Server) handleConnection(conn net.Conn) error {
	bufferLen := 8
	buffer := make([]byte, bufferLen)
	p := http.NewParser()
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return err
		}
		if err := p.Feed(buffer[:n]); err != nil {
			return err
		}
	}
}
