package server

import (
	"fmt"
	"net"

	"github.com/tusharameria/go-http/internal/http"
)

func (s *Server) handleConnection(conn net.Conn) error {
	reader := http.NewReader(conn)
	request, err := reader.ReadRequest()
	if err != nil {
		return err
	}
	fmt.Println(request)
	return nil
}
