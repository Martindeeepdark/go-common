package otel

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestServerOptions_ReturnsOptions(t *testing.T) {
	opts := ServerOptions()
	require.NotEmpty(t, opts, "ServerOptions should return at least one option")
}

func TestServerOptions_NonNil(t *testing.T) {
	opts := ServerOptions()
	for _, opt := range opts {
		assert.NotNil(t, opt)
	}
}

func TestServerOptions_Type(t *testing.T) {
	opts := ServerOptions()
	assert.IsType(t, []grpc.ServerOption{}, opts)
}

func TestClientOptions_ReturnsOptions(t *testing.T) {
	opts := ClientOptions()
	require.NotEmpty(t, opts, "ClientOptions should return at least one option")
}

func TestClientOptions_NonNil(t *testing.T) {
	opts := ClientOptions()
	for _, opt := range opts {
		assert.NotNil(t, opt)
	}
}

func TestClientOptions_Type(t *testing.T) {
	opts := ClientOptions()
	assert.IsType(t, []grpc.DialOption{}, opts)
}
