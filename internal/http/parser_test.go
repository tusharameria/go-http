package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestParseRequestLine(t *testing.T) {
	raw := []byte("GET / HTTP/1.1")
	p := &Parser{}
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
