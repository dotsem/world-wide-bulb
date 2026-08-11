// Package utils provides utility functions and helper services.
package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hasher provides salted cryptographic hashing for IP addresses.
type Hasher struct {
	salt string
}

// NewHasher creates a new Hasher initialized with the given salt.
func NewHasher(salt string) *Hasher {
	return &Hasher{salt: salt}
}

// Hash returns a hex-encoded SHA-256 hash of the IP combined with the salt.
func (h *Hasher) Hash(ip string) string {
	hash := sha256.Sum256([]byte(ip + h.salt))
	return hex.EncodeToString(hash[:])
}
