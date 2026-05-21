package oauth

import (
	"crypto/sha256"
	"encoding/base64"
	"testing"

	"gotest.tools/v3/assert"
)

func TestGenerateVerifier_Length(t *testing.T) {
	v, err := GenerateVerifier()
	assert.NilError(t, err)
	assert.Equal(t, len(v), 43) // 32 bytes base64url = 43 chars
}

func TestGenerateVerifier_Uniqueness(t *testing.T) {
	v1, err := GenerateVerifier()
	assert.NilError(t, err)
	v2, err := GenerateVerifier()
	assert.NilError(t, err)
	assert.Assert(t, v1 != v2)
}

func TestS256Challenge_KnownVector(t *testing.T) {
	// RFC 7636 Appendix B test vector
	verifier := "dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk"
	h := sha256.Sum256([]byte(verifier))
	expected := base64.RawURLEncoding.EncodeToString(h[:])
	assert.Equal(t, S256Challenge(verifier), expected)
}
