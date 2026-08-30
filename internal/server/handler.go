package server

import (
	"net"

	"github.com/tusharameria/go-http/internal/http"
)

func (s *Server) handleConnection(conn net.Conn) error {
	defer conn.Close()

	connection := http.NewConnection(conn)
	return connection.Serve()
}
