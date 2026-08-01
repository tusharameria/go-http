package http

import (
	"bytes"
	"fmt"
)

func Parse(data []byte) (*Request, error) {
	line, _, _ := bytes.Cut(data, []byte("\r\n"))
	parts := bytes.Split(line, []byte(" "))
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid request line")
	}

	return &Request{
		Method:  string(parts[0]),
		Path:    string(parts[1]),
		Version: string(parts[2]),
	}, nil
}
