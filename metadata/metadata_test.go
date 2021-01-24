package metadata

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClientMetadata(t *testing.T) {
	ctx := context.Background()
	md := ClientMetadata(ctx)
	assert.NotNil(t,md)

	md["test"] = []byte("test_client_metadata")
	newCtx := WithClientMetadata(ctx, md)
	md = ClientMetadata(newCtx)
	assert.Equal(t, string(md["test"]), "test_client_metadata")
}

func TestServerMetadata(t *testing.T) {
	ctx := context.Background()
	md := ServerMetadata(ctx)
	assert.NotNil(t,md)

	md["test"] = []byte("test_server_metadata")
	newCtx := WithServerMetadata(ctx, md)
	md1 := ServerMetadata(newCtx)
	assert.Equal(t, string(md1["test"]), "test_server_metadata")
}