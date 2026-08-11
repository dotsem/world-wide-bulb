package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

type Hasher struct {
	salt string
}

func NewHasher(salt string) *Hasher {
	return &Hasher{salt: salt}
}

func (h *Hasher) Hash(ip string) string {
	hash := sha256.Sum256([]byte(ip + h.salt))
	return hex.EncodeToString(hash[:])
}
