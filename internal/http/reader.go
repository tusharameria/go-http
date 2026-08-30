package http

import (
	"io"
	"net"
)

type Reader struct {
	conn   net.Conn
	parser *Parser
	buffer []byte
}

func NewReader(conn net.Conn) *Reader {
	bufferLen := 8
	buffer := make([]byte, bufferLen)

	return &Reader{
		conn:   conn,
		parser: NewParser(),
		buffer: buffer,
	}
}

func (r *Reader) ReadRequest() (*Request, error) {
	for {
		n, err := r.conn.Read(r.buffer)

		if n > 0 {
			request, parseErr := r.parser.Feed(r.buffer[:n])
			if parseErr != nil {
				return nil, parseErr
			}

			if request != nil {
				return request, nil
			}
		}

		if err != nil {
			if err == io.EOF {
				return nil, nil
			}
			return nil, err
		}
	}
}
