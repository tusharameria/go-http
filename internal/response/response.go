package response

import (
	"fmt"
	"io"

	"github.com/tusharameria/go-http/internal/headers"
)

type StatusCode int

const (
	StatusOK                  StatusCode = 200
	StatusBadRequest          StatusCode = 400
	StatusInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	str := "Status Not Handled"
	switch statusCode {
	case StatusOK:
		str = "HTTP/1.1 200 OK"
	case StatusBadRequest:
		str = "HTTP/1.1 400 Bad Request"
	case StatusInternalServerError:
		str = "HTTP/1.1 500 Internal Server Error"
	default:
		return fmt.Errorf(str)
	}
	str += "\r\n"
	_, err := w.Write([]byte(str))
	return err
}

func GetDefaultHeaders(contentLen int) *headers.Headers {
	h := headers.NewHeaders()
	h.Set("Content-Length", fmt.Sprintf("%d", contentLen))
	h.Set("Connection", "close")
	h.Set("Content-Type", "text/plain")
	return h
}

func WriteHeaders(w io.Writer, headers *headers.Headers) error {
	var err error
	headers.ForEach(func(k, v string) {
		_, err = w.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return
		}
	})
	w.Write([]byte("\r\n"))
	return err
}
