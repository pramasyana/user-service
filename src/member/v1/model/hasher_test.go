package model

import (
	"crypto/sha1"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	// TestPassPlain const
	TestPassPlain = "wannabenaked"
	// TestPassSalt const
	TestPassSalt = "15000.Efns3HvrXL7aP75bjMpzlwWp7s/yZwlEdRNuyuTWZk9JGXXCyGtBa/BQTz3s4vY8ewgDcgL/xz2efTYzPdrKhg=="
	// TestPassCipher const
	TestPassCipher = "5fJx5lYkAtHGVFktXjvyAL8t9ZOyjjMGhXXy7QohBUFwYHG9DZKWdhPd1LE/tA9nZcbGHUryKNP+VvH99TLG3w=="
)

func TestBMDHash(t *testing.T) {
	hasher := NewPBKDF2Hasher(SaltSize, SaltSize, IterationsCount, sha1.New) // use 15001 to offset to 1

	t.Run("Testing Hash Verification", func(t *testing.T) {
		assert.Equal(t, nil, hasher.ParseSalt(TestPassPlain))
		assert.NoError(t, hasher.ParseSalt(TestPassSalt))
		assert.Equal(t, IterationsCount, hasher.iteration)
		assert.Equal(t, TestPassCipher, base64.StdEncoding.EncodeToString(hasher.Hash([]byte(TestPassPlain))))
	})

}

func TestGenerateSalt(t *testing.T) {
	hasher := NewPBKDF2Hasher(SaltSize, SaltSize, IterationsCount, sha1.New) // use 15001 to offset to 1

	t.Run("Testing Generate Salt", func(t *testing.T) {
		generated := hasher.GenerateSalt()
		assert.Equal(t, generated, generated)
	})
}
