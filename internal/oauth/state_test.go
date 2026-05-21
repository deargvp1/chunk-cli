package oauth

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestGenerateState_Length(t *testing.T) {
	s, err := GenerateState()
	assert.NilError(t, err)
	assert.Equal(t, len(s), 32) // 16 bytes hex = 32 chars
}

func TestGenerateState_Uniqueness(t *testing.T) {
	s1, err := GenerateState()
	assert.NilError(t, err)
	s2, err := GenerateState()
	assert.NilError(t, err)
	assert.Assert(t, s1 != s2)
}
