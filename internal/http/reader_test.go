package http

import (
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReaderReadRequest(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	reader := NewReader(serverConn)

	go func() {
		_, _ = clientConn.Write([]byte(
			"GET /hello HTTP/1.1\r\n" +
				"Host: localhost\r\n" +
				"\r\n",
		))
		_ = clientConn.Close()
	}()

	request, err := reader.ReadRequest()

	require.NoError(t, err)
	require.NotNil(t, request)
	require.Equal(t, "GET", request.requestLine.Method)
	require.Equal(t, "/hello", request.requestLine.Path)
	require.Equal(t, "HTTP/1.1", request.requestLine.Version)
	require.Equal(t, "localhost", request.headers.Get("Host"))
	require.Empty(t, request.body)
}

func TestReaderReadRequest_Fragmented(t *testing.T) {
	serverConn, clientConn := net.Pipe()
	defer serverConn.Close()
	defer clientConn.Close()

	reader := NewReader(serverConn)

	go func() {
		parts := []string{
			"POST /hello HTTP/1.1\r\n",
			"Host: localhost\r\n",
			"Content-Length: 11\r\n",
			"\r\n",
			"Hello World",
		}

		for _, part := range parts {
			_, _ = clientConn.Write([]byte(part))
		}

		_ = clientConn.Close()
	}()

	request, err := reader.ReadRequest()

	require.NoError(t, err)
	require.NotNil(t, request)
	require.Equal(t, "POST", request.requestLine.Method)
	require.Equal(t, "/hello", request.requestLine.Path)
	require.Equal(t, "localhost", request.headers.Get("Host"))
	require.Equal(t, []byte("Hello World"), request.body)
}
