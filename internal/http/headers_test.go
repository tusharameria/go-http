package http

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHeaders_GetReturnsFirstValue(t *testing.T) {
	headers := NewHeaders()

	headers.add("Host", "localhost")
	headers.add("Host", "127.0.0.1")

	require.Equal(t, "localhost", headers.Get("Host"))
}

func TestHeaders_GetIsCaseInsensitive(t *testing.T) {
	headers := NewHeaders()

	headers.add("Host", "localhost")

	require.Equal(t, "localhost", headers.Get("Host"))
	require.Equal(t, "localhost", headers.Get("HOST"))
	require.Equal(t, "localhost", headers.Get("host"))
	require.Equal(t, "localhost", headers.Get("HoSt"))
}

func TestHeaders_ValuesReturnsAllValues(t *testing.T) {
	headers := NewHeaders()

	headers.add("Accept", "application/json")
	headers.add("Accept", "application/xml")

	values := headers.Values("Accept")

	require.Len(t, values, 2)
	require.Equal(t, []string{
		"application/json",
		"application/xml",
	}, values)
}

func TestHeaders_ValuesReturnsCopy(t *testing.T) {
	headers := NewHeaders()

	headers.add("Accept", "application/json")

	values := headers.Values("Accept")
	values[0] = "modified"

	require.Equal(t, "application/json", headers.Get("Accept"))
	require.Equal(t, []string{"application/json"}, headers.Values("Accept"))
}

func TestHeaders_Has(t *testing.T) {
	headers := NewHeaders()

	headers.add("Host", "localhost")

	require.True(t, headers.Has("Host"))
	require.True(t, headers.Has("HOST"))
	require.True(t, headers.Has("host"))

	require.False(t, headers.Has("Content-Type"))
}

func TestHeaders_GetMissingHeader(t *testing.T) {
	headers := NewHeaders()

	require.Equal(t, "", headers.Get("Missing"))
}

func TestHeaders_ValuesMissingHeader(t *testing.T) {
	headers := NewHeaders()

	require.Nil(t, headers.Values("Missing"))
}
