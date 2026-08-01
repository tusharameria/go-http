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
	require.Equal(t, "localhost", req.Headers["Host"])
	require.Equal(t, "curl", req.Headers["User-Agent"])
}
