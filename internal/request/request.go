package request

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"

	"github.com/tusharameria/go-http/internal/headers"
)

type parserState string

const (
	StateInit    parserState = "init"
	StateHeaders parserState = "headers"
	StateBody    parserState = "body"
	StateDone    parserState = "done"
	StateError   parserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	Headers     *headers.Headers
	state       parserState
	Body        string
}

func getInt(headers *headers.Headers, name string, defaultValue int) int {
	valueStr, exists := headers.Get(name)
	if !exists {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

func newRequest() *Request {
	return &Request{
		state:   StateInit,
		Headers: headers.NewHeaders(),
	}
}

var ErrorMalformedRequestLine = fmt.Errorf("malformed request-line")
var ErrorUnsupportedHttpVersion = fmt.Errorf("unsupported http version")
var ErrorMalformedContentLength = fmt.Errorf("malformed content-length")
var ErrorContentLengthMismatch = fmt.Errorf("content length mismatch")
var ErrorRequestInErrorState = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n")
var ContentLength = "Content-Length"

// func parseBody(b []byte, contentLength int) ([]byte, int, error) {

// }

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	startLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(startLine, []byte(" "))
	if len(parts) != 3 {
		return nil, 0, ErrorMalformedRequestLine
	}

	httpVersionParts := bytes.Split(parts[2], []byte("/"))
	if len(httpVersionParts) != 2 ||
		string(httpVersionParts[0]) != "HTTP" ||
		string(httpVersionParts[1]) != "1.1" {
		return nil, 0, ErrorUnsupportedHttpVersion
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpVersionParts[1]),
	}

	return rl, read, nil
}

func (r *Request) hasBody() bool {
	// TODO : Implement chunk
	length := getInt(r.Headers, ContentLength, 0)
	return length != 0
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0
dance:
	for {
		currentData := data[read:]
		if len(currentData) == 0 {
			break dance
		}
		switch r.state {
		case StateError:
			return 0, ErrorRequestInErrorState
		case StateInit:
			rl, n, err := parseRequestLine(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break dance
			}
			r.RequestLine = *rl
			read += n
			r.state = StateHeaders
		case StateHeaders:
			n, done, err := r.Headers.Parse(currentData)
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break dance
			}
			read += n
			if done {
				if r.hasBody() {
					r.state = StateBody
				} else {
					r.state = StateDone
				}
			}
		case StateBody:
			length := getInt(r.Headers, ContentLength, 0)
			if length == 0 {
				panic("chunk not implemented")
			}

			remaining := min(length-len(r.Body), len(currentData))
			r.Body += string(currentData[:remaining])
			read += remaining

			if len(r.Body) == length {
				r.state = StateDone
			}
		case StateDone:
			break dance
		default:
			panic("bad programmer!!!")
		}
	}
	return read, nil
}

func (r *Request) done() bool {
	return r.state == StateDone
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()
	buff := make([]byte, 1024)
	buffLen := 0
	for !request.done() {
		n, err := reader.Read(buff[buffLen:])
		if err != nil {
			return nil, errors.Join(
				fmt.Errorf("unable to reader.Read"),
				err,
			)
		}
		buffLen += n
		readN, newErr := request.parse(buff[:buffLen])
		if newErr != nil {
			return nil, newErr
		}
		copy(buff, buff[readN:buffLen])
		buffLen -= readN

		// if err == io.EOF {
		// 	fmt.Println("EOF")
		// 	break
		// }
	}

	return request, nil

}
