package server

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewWithListener(t *testing.T) {
	srv := NewWithListener(nil)

	require.NotNil(t, srv)
	require.Nil(t, srv.listener)
}
