package oauth

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestListenForCallback_Success(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	port, resultCh, cleanup, err := ListenForCallback(ctx)
	assert.NilError(t, err)
	defer cleanup()

	url := fmt.Sprintf("http://127.0.0.1:%d/callback?code=test-code&state=test-state", port)
	resp, err := http.Get(url)
	assert.NilError(t, err)
	resp.Body.Close()
	assert.Equal(t, resp.StatusCode, http.StatusOK)

	res := <-resultCh
	assert.Equal(t, res.Code, "test-code")
	assert.Equal(t, res.State, "test-state")
	assert.Equal(t, res.Error, "")
}

func TestListenForCallback_Error(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	port, resultCh, cleanup, err := ListenForCallback(ctx)
	assert.NilError(t, err)
	defer cleanup()

	url := fmt.Sprintf("http://127.0.0.1:%d/callback?error=access_denied&state=test-state", port)
	resp, err := http.Get(url)
	assert.NilError(t, err)
	resp.Body.Close()

	res := <-resultCh
	assert.Equal(t, res.Error, "access_denied")
	assert.Equal(t, res.State, "test-state")
	assert.Equal(t, res.Code, "")
}

func TestListenForCallback_ContextCancel(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	_, _, cleanup, err := ListenForCallback(ctx)
	assert.NilError(t, err)
	defer cleanup()

	cancel()
	// Server should shut down without hanging; cleanup is the verification.
}
