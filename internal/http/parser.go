package http

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type parserState int
type bodyType int

const contentLengthKey = "Content-Length"
const crlf = "\r\n"
const (
	stateRequestLine parserState = iota
	stateHeaders
	stateBody
	stateComplete
)
const (
	bodyNone bodyType = iota
	bodyContentLength
	bodyChunked
)

var crlfBytes = []byte(crlf)

type Parser struct {
	state        parserState
	requestLine  RequestLine
	headers      *Headers
	body         []byte
	bodyType     bodyType
	bodyMetaData int
	buf          []byte
}

func NewParser() *Parser {
	return &Parser{
		state:    stateRequestLine,
		headers:  NewHeaders(),
		bodyType: bodyNone,
	}
}

func (p *Parser) Feed(data []byte) (*Request, error) {
	p.buf = append(p.buf, data...)

	for {
		switch p.state {
		case stateRequestLine:
			idx := bytes.Index(p.buf, crlfBytes)
			if idx == -1 {
				return nil, nil
			}

			requestLine, err := p.parseRequestLine(p.buf[:idx])
			if err != nil {
				return nil, err
			}

			p.requestLine = requestLine
			p.buf = p.buf[idx+len(crlfBytes):]
			p.state = stateHeaders

		case stateHeaders:
			idx := bytes.Index(p.buf, crlfBytes)
			if idx == -1 {
				return nil, nil
			}

			if idx == 0 {
				if err := p.determineBody(); err != nil {
					return nil, err
				}

				p.buf = p.buf[len(crlfBytes):]
				p.state = stateBody
				continue
			}

			if err := p.parseHeaderLine(p.buf[:idx]); err != nil {
				return nil, err
			}

			p.buf = p.buf[idx+len(crlfBytes):]

		case stateBody:
			switch p.bodyType {
			case bodyNone:
				p.state = stateComplete

			case bodyContentLength:
				if len(p.buf) == 0 {
					return nil, nil
				}

				idx := min(p.bodyMetaData, len(p.buf))

				p.body = append(p.body, p.buf[:idx]...)
				p.buf = p.buf[idx:]
				p.bodyMetaData -= idx

				if p.bodyMetaData > 0 {
					return nil, nil
				}

				p.state = stateComplete
			}

		case stateComplete:
			return &Request{
				requestLine: p.requestLine,
				headers:     p.headers,
				body:        p.body,
			}, nil
		}
	}
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

func (p *Parser) parseHeaderLine(line []byte) error {
	headerLineParts := bytes.SplitN(line, []byte(":"), 2)
	if len(headerLineParts) != 2 {
		return fmt.Errorf("invalid header line")
	}
	key := strings.TrimSpace(string(headerLineParts[0]))
	value := strings.TrimSpace(string(headerLineParts[1]))
	p.headers.add(key, value)

	return nil
}

func (p *Parser) determineBody() error {
	if !p.headers.Has(contentLengthKey) {
		p.bodyType = bodyNone
		return nil
	}

	contentLengthValue, err := strconv.Atoi(p.headers.Get(contentLengthKey))
	if err != nil {
		return err
	}
	p.bodyType = bodyContentLength
	p.bodyMetaData = contentLengthValue

	return nil
}
