package server

import "net"

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
	conn, err := s.listener.Accept()
	if err != nil {
		return err
	}
	defer conn.Close()

	return s.handleConnection(conn)
}
