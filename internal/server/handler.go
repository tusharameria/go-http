package server

import (
	"bytes"
	"fmt"
	"net"
)

func (s *Server) handleConnection(conn net.Conn) error {
	bufferLen := 8
	buffer := make([]byte, bufferLen)
	var aggregator []byte
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			return err
		}
		aggregator = append(aggregator, buffer[:n]...)
		lastIdx := bytes.Index(aggregator, []byte("\n"))
		if lastIdx >= 0 {
			fmt.Printf("Full Text : %s\n", aggregator[:lastIdx])
			conn.Write([]byte(fmt.Sprintf("Response from server : Hey you sent us \"%s\"\n", aggregator[:lastIdx])))
			aggregator = aggregator[lastIdx+1:]
		}
	}
}
