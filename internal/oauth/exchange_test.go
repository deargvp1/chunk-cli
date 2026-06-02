package oauth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gotest.tools/v3/assert"
)

func TestExchangeToken_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, http.MethodPost)
		assert.Equal(t, r.URL.Path, "/oauth/token")
		assert.Equal(t, r.Header.Get("Content-Type"), "application/x-www-form-urlencoded")

		assert.NilError(t, r.ParseForm())
		assert.Equal(t, r.FormValue("grant_type"), "authorization_code")
		assert.Equal(t, r.FormValue("code"), "test-code")
		assert.Equal(t, r.FormValue("client_id"), ClientID)
		assert.Equal(t, r.FormValue("code_verifier"), "test-verifier")

		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"access_token":"CCIPAT_test_token","token_type":"Bearer","expires_in":7776000}`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	tok, err := exchangeToken(ctx, srv.URL, "test-code", "http://127.0.0.1:12345/callback", "test-verifier")
	assert.NilError(t, err)
	assert.Equal(t, tok.AccessToken, "CCIPAT_test_token")
	assert.Equal(t, tok.TokenType, "Bearer")
	assert.Equal(t, tok.ExpiresIn, 7776000)
}

func TestExchangeToken_ErrorResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(`{"error":"invalid_grant"}`))
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := exchangeToken(ctx, srv.URL, "bad-code", "http://127.0.0.1:12345/callback", "test-verifier")
	assert.ErrorContains(t, err, "400")
}
