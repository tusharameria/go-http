package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRequestLine(t *testing.T) {
	raw := []byte("GET / HTTP/1.1")
	p := &Parser{
		state: stateRequestLine,
	}
	req, err := p.parseRequestLine(raw)

	require.NoError(t, err)
	require.Equal(t, "GET", req.Method)
	require.Equal(t, "/", req.Path)
	require.Equal(t, "HTTP/1.1", req.Version)
}

func TestParseRequestLine_Invalid(t *testing.T) {
	raw := []byte("GET /\r\n")
	p := &Parser{}
	req, err := p.Parse(raw)

	require.Error(t, err)
	require.Nil(t, req)
}

func TestParseRequest_WithHeaders(t *testing.T) {
	raw := []byte(
		"GET / HTTP/1.1\r\n" +
			"Host: localhost\r\n" +
			"User-Agent: curl\r\n" +
			"\r\n",
	)

	p := &Parser{}
	req, err := p.Parse(raw)

	require.NoError(t, err)
	require.Equal(t, 1, len(req.Headers.Values("Host")))
	require.Equal(t, true, req.Headers.Has("Host"))
	require.Equal(t, true, req.Headers.Has("host"))
	require.Equal(t, "localhost", req.Headers.Get("Host"))
	require.Equal(t, 1, len(req.Headers.Values("User-Agent")))
	require.Equal(t, true, req.Headers.Has("User-Agent"))
	require.Equal(t, true, req.Headers.Has("UsEr-AgeNt"))
	require.Equal(t, "curl", req.Headers.Get("User-Agent"))
}

func TestParserFeed_RequestLineComplete(t *testing.T) {
	p := &Parser{
		state: stateRequestLine,
	}

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
	p := &Parser{
		state: stateRequestLine,
	}

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
