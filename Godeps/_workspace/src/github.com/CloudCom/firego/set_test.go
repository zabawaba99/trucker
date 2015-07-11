package firego

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSet(t *testing.T) {
	t.Parallel()
	var (
		response = `{"foo":"bar"}`
		server   = newTestServer(response)
		fb       = New(server.URL)
	)
	defer server.Close()

	fb.Set(response)
	require.Len(t, server.receivedReqs, 1)

	req := server.receivedReqs[0]
	assert.Equal(t, "PUT", req.Method)
}
