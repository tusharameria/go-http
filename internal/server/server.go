package server

import (
	"fmt"
	"net"
)

type Server struct {
	listener net.Listener
}

func New(addr string) (*Server, error) {
	srv, err := net.Listen("tcp", ":8082")
	if err != nil {
		return nil, err
	}
	return &Server{
		listener: srv,
	}, nil
}

func NewWithListener(listener net.Listener) *Server {
	return &Server{
		listener: listener,
	}
}

func (s *Server) Serve() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			return err
		}

		go func() {
			if err := s.handleConnection(conn); err != nil {
				fmt.Println(err)
			}
		}()
	}
}
