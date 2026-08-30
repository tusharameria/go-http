package http

import (
	"fmt"
	"net"
)

type Connection struct {
	conn   net.Conn
	reader *Reader
}

func NewConnection(conn net.Conn) *Connection {
	return &Connection{
		conn:   conn,
		reader: NewReader(conn),
	}
}

func (c *Connection) Serve() error {
	request, err := c.reader.ReadRequest()
	if err != nil {
		return err
	}
	fmt.Println(request)
	return nil
}
