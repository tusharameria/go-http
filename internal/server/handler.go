package server

import (
	"fmt"
	"io"
	"net"

	"github.com/tusharameria/go-http/internal/http"
)

func (s *Server) handleConnection(conn net.Conn) error {
	bufferLen := 8
	buffer := make([]byte, bufferLen)
	p := http.NewParser()
	for {
		n, err := conn.Read(buffer)

		if n > 0 {
			request, parseErr := p.Feed(buffer[:n])
			if parseErr != nil {
				return parseErr
			}

			if request != nil {
				// TODO : handle request
				fmt.Println(request)
				return nil
			}
		}

		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
