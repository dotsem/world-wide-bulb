package utils_test

import (
	"testing"
	"world-wide-bulb/internal/utils"

	"github.com/stretchr/testify/assert"
)

func TestHasher(t *testing.T) {
	t.Run("produces deterministic hash for same input and salt", func(t *testing.T) {
		hasher := utils.NewHasher("test_salt")

		hash1 := hasher.Hash("192.168.1.1")
		hash2 := hasher.Hash("192.168.1.1")

		assert.NotEmpty(t, hash1)
		assert.Equal(t, hash1, hash2)
	})

	t.Run("different IPs produce distinct hashes", func(t *testing.T) {
		hasher := utils.NewHasher("test_salt")

		hash1 := hasher.Hash("192.168.1.1")
		hash2 := hasher.Hash("192.168.1.2")

		assert.NotEqual(t, hash1, hash2)
	})

	t.Run("different salts produce distinct hashes for same IP", func(t *testing.T) {
		hasher1 := utils.NewHasher("salt_a")
		hasher2 := utils.NewHasher("salt_b")

		hash1 := hasher1.Hash("192.168.1.1")
		hash2 := hasher2.Hash("192.168.1.1")

		assert.NotEqual(t, hash1, hash2)
	})
}
