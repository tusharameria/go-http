package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
)

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

type Request struct {
	RequestLine RequestLine
	State       parserState
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

var ERROR_MALFORMED_REQUEST_LINE = fmt.Errorf("malformed request-line")
var ERROR_UNSUPPORTED_HTTP_VERSION = fmt.Errorf("unsupported http version")
var SEPARATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ERROR_MALFORMED_REQUEST_LINE
	}

	httpVersionParts := bytes.Split(parts[2], []byte("/"))
	if len(httpVersionParts) != 2 ||
		string(httpVersionParts[0]) != "HTTP" ||
		string(httpVersionParts[1]) != "1.1" {
		return nil, 0, ERROR_UNSUPPORTED_HTTP_VERSION
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpVersionParts[1]),
	}

	return rl, read, nil
}

func newRequest() *Request {
	return &Request{
		State: StateInit,
	}
}

func (r *Request) parse(b []byte) (int, error) {
	rl, readN, err := parseRequestLine(b)
	if err != nil {
		return 0, err
	}
	if readN != 0 {
		r.RequestLine = *rl
		r.State = StateDone
	}
	return 0, nil
}

func (r *Request) done() bool {
	return r.State == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buff := make([]byte, 1024)
	buffLen := 0
	for !request.done() {
		n, err := reader.Read(buff[buffLen:])
		if err != nil && err != io.EOF {
			return nil, errors.Join(
				fmt.Errorf("unable to reader.Read"),
				err,
			)
		}

		buffLen += n
		readN, newErr := request.parse(buff[:buffLen])
		if newErr != nil {
			return nil, err
		}
		copy(buff, buff[readN:buffLen])
		buffLen -= readN

		if err == io.EOF {
			break
		}
	}

	return request, nil

}
