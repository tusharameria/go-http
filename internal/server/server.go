package server

import (
	"fmt"
	"io"
	"net"
	"sync/atomic"

	"github.com/tusharameria/go-http/internal/response"
)

type Server struct {
	port     uint16
	open     atomic.Bool
	listener net.Listener
}

func runConnection(s *Server, conn io.ReadWriteCloser) {
	defer conn.Close()
	response.WriteStatusLine(conn, response.StatusOK)
	h := response.GetDefaultHeaders(13)
	response.WriteHeaders(conn, h)
}

func runServer(s *Server, listener net.Listener) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			return
		}
		if s.open.Load() == false {
			return
		}
		go runConnection(s, conn)
	}
}

func Serve(port uint16) (*Server, error) {
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}
	s := &Server{
		port: port,
	}
	s.open.Store(true)
	go runServer(s, listener)
	return s, nil
}

func (s *Server) Close() error {
	if err := s.listener.Close(); err != nil {
		return err
	}
	s.open.Store(false)
	return nil
}
