package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRequestLine(t *testing.T) {
	raw := []byte("GET / HTTP/1.1")
	p := NewParser()
	req, err := p.parseRequestLine(raw)

	require.NoError(t, err)
	require.Equal(t, "GET", req.Method)
	require.Equal(t, "/", req.Path)
	require.Equal(t, "HTTP/1.1", req.Version)
}

func TestParserFeed_RequestLineComplete(t *testing.T) {
	p := NewParser()

	err := p.Feed([]byte("GET /hello HTTP/1.1\r\n"))

	require.NoError(t, err)
	require.Equal(t, stateHeaders, p.state)
	require.Equal(t, RequestLine{
		Method:  "GET",
		Path:    "/hello",
		Version: "HTTP/1.1",
	}, p.requestLine)
	require.Empty(t, p.buf)
}

func TestParserFeed_RequestLineSplitAcrossFeeds(t *testing.T) {
	p := NewParser()

	err := p.Feed([]byte("GET /hel"))
	require.NoError(t, err)
	require.Equal(t, stateRequestLine, p.state)
	require.Empty(t, p.requestLine)
	require.Equal(t, []byte("GET /hel"), p.buf)

	err = p.Feed([]byte("lo HTTP/1.1\r\n"))
	require.NoError(t, err)
	require.Equal(t, stateHeaders, p.state)
	require.Equal(t, RequestLine{
		Method:  "GET",
		Path:    "/hello",
		Version: "HTTP/1.1",
	}, p.requestLine)
	require.Empty(t, p.buf)
}

func TestParserFeed_HeadersComplete(t *testing.T) {
	p := NewParser()

	err := p.Feed([]byte(
		"GET / HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"Content-Type: application/json\r\n" +
			"\r\n",
	))

	require.NoError(t, err)
	require.Equal(t, stateBody, p.state)

	require.Equal(t, "localhost", p.headers.Get("Host"))
	require.Equal(t, "application/json", p.headers.Get("Content-Type"))
	require.Empty(t, p.buf)
}

func TestParserFeed_HeaderSplitAcrossFeeds(t *testing.T) {
	p := NewParser()

	err := p.Feed([]byte(
		"GET / HTTP/1.1\r\n" +
			"Host: local",
	))
	require.NoError(t, err)
	require.Equal(t, stateHeaders, p.state)
	require.Equal(t, []byte("Host: local"), p.buf)

	err = p.Feed([]byte(
		"host\r\n" +
			"Content-Type: application/json\r\n" +
			"\r\n",
	))
	require.NoError(t, err)
	require.Equal(t, stateBody, p.state)

	require.Equal(t, "localhost", p.headers.Get("Host"))
	require.Equal(t, "application/json", p.headers.Get("Content-Type"))
	require.Empty(t, p.buf)
}

func TestParserFeed_MultipleHeadersInSingleFeed(t *testing.T) {
	p := NewParser()

	err := p.Feed([]byte(
		"GET / HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"Accept: application/json\r\n" +
			"Connection: keep-alive\r\n" +
			"\r\n",
	))

	require.NoError(t, err)
	require.Equal(t, stateBody, p.state)

	require.Equal(t, "localhost", p.headers.Get("Host"))
	require.Equal(t, "application/json", p.headers.Get("Accept"))
	require.Equal(t, "keep-alive", p.headers.Get("Connection"))
	require.Empty(t, p.buf)
}

func TestParserFeed_NoBody(t *testing.T) {
	p := NewParser()

	err := p.Feed([]byte(
		"GET / HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"\r\n",
	))

	require.NoError(t, err)
	require.Equal(t, bodyNone, p.bodyType)
	require.Zero(t, p.bodyMetaData)
	require.Equal(t, stateBody, p.state)
}

func TestParserFeed_ContentLength(t *testing.T) {
	p := NewParser()

	err := p.Feed([]byte(
		"POST / HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"Content-Length: 42\r\n" +
			"\r\n",
	))

	require.NoError(t, err)
	require.Equal(t, bodyContentLength, p.bodyType)
	require.Equal(t, 42, p.bodyMetaData)
	require.Equal(t, stateBody, p.state)
}

func TestParserFeed_ContentLengthBody(t *testing.T) {
	p := NewParser()

	err := p.Feed([]byte(
		"POST / HTTP/1.1\r\n" +
			"Content-Length: 11\r\n" +
			"\r\n" +
			"Hello World",
	))

	require.NoError(t, err)
	require.Equal(t, stateComplete, p.state)
	require.Equal(t, []byte("Hello World"), p.body)
	require.Empty(t, p.buf)
	require.Zero(t, p.bodyMetaData)
}

func TestParserFeed_ContentLengthBodySplitAcrossFeeds(t *testing.T) {
	p := NewParser()

	err := p.Feed([]byte(
		"POST / HTTP/1.1\r\n" +
			"Content-Length: 11\r\n" +
			"\r\n" +
			"Hello",
	))

	require.NoError(t, err)
	require.Equal(t, stateBody, p.state)
	require.Equal(t, []byte("Hello"), p.body)
	require.Equal(t, 6, p.bodyMetaData)

	err = p.Feed([]byte(" World"))

	require.NoError(t, err)
	require.Equal(t, stateComplete, p.state)
	require.Equal(t, []byte("Hello World"), p.body)
	require.Empty(t, p.buf)
	require.Zero(t, p.bodyMetaData)
}

func TestParserFeed_ContentLengthBodyAndExtraData(t *testing.T) {
	p := NewParser()

	err := p.Feed([]byte(
		"POST / HTTP/1.1\r\n" +
			"Content-Length: 5\r\n" +
			"\r\n" +
			"HelloExtra",
	))

	require.NoError(t, err)
	require.Equal(t, stateComplete, p.state)
	require.Equal(t, []byte("Hello"), p.body)
	require.Equal(t, []byte("Extra"), p.buf)
}
