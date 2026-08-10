package http

import (
	"bytes"
	"fmt"
	"strings"
)

type parserState int

const crlf = "\r\n"
const (
	stateRequestLine parserState = iota
	stateHeaders
)

var crlfBytes = []byte(crlf)

type Parser struct {
	state       parserState
	requestLine RequestLine
	buf         []byte
}

func (p *Parser) Feed(data []byte) error {
	p.buf = append(p.buf, data...)
	idx := bytes.Index(p.buf, crlfBytes)
	if idx == -1 {
		return nil
	}
	requestLine, err := p.parseRequestLine(p.buf[:idx])
	if err != nil {
		return err
	}
	p.requestLine = requestLine
	p.buf = p.buf[idx+len(crlfBytes):]
	p.state = stateHeaders
	return nil
}

func (p *Parser) parseRequestLine(line []byte) (RequestLine, error) {
	requestLine := RequestLine{}
	rlParts := bytes.Split(line, []byte(" "))
	if len(rlParts) != 3 {
		return requestLine, fmt.Errorf("invalid request line")
	}
	requestLine.Method = string(rlParts[0])
	requestLine.Path = string(rlParts[1])
	requestLine.Version = string(rlParts[2])

	return requestLine, nil
}

func (p *Parser) parseHeaders(line []byte) (*Headers, error) {
	headers := NewHeaders()
	if len(line) != 0 {
		headerParts := bytes.Split(line, []byte("\r\n"))
		lenHeadParts := len(headerParts)
		if len(headerParts[lenHeadParts-1]) != 0 || len(headerParts[lenHeadParts-2]) != 0 {
			return nil, fmt.Errorf("headers not ended properly")
		}

		for i := 0; i < lenHeadParts-2; i++ {
			headerLine := headerParts[i]
			headerLineParts := bytes.SplitN(headerLine, []byte(":"), 2)
			if len(headerLineParts) != 2 {
				return nil, fmt.Errorf("invalid header line")
			}
			key := strings.TrimSpace(string(headerLineParts[0]))
			value := strings.TrimSpace(string(headerLineParts[1]))
			headers.add(key, value)
		}
	}

	return headers, nil
}

func (p *Parser) Parse(data []byte) (*Request, error) {
	requesLineBytes, restOfData, found := bytes.Cut(data, []byte("\r\n"))
	if !found {
		return nil, fmt.Errorf("request line not ended")
	}

	requestLine, err := p.parseRequestLine(requesLineBytes)
	if err != nil {
		return nil, err
	}

	headers, err := p.parseHeaders(restOfData)
	if err != nil {
		return nil, err
	}

	return &Request{
		RequestLine: requestLine,
		Headers:     headers,
	}, nil
}
