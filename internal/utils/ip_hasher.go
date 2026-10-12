// Package utils provides utility functions and helper services.
package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

// Hasher provides salted cryptographic hashing and HMAC signing for IP addresses and cookies.
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

// Sign creates an HMAC-SHA256 signature for data using the salt as secret key.
func (h *Hasher) Sign(data string) string {
	mac := hmac.New(sha256.New, []byte(h.salt))
	mac.Write([]byte(data))
	return hex.EncodeToString(mac.Sum(nil))
}

// Verify checks whether sig is a valid HMAC-SHA256 signature for data.
func (h *Hasher) Verify(data, sig string) bool {
	expected := h.Sign(data)
	return hmac.Equal([]byte(expected), []byte(sig))
}
