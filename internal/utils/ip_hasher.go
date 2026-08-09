package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

var salt string

func init() {
	var err error
	salt, err = GetEnvOrErrInProd("IP_SALT", "dev_fallback")
	if err != nil {
		panic(err)
	}
}

func HashIP(ip string) string {
	hash := sha256.Sum256([]byte(ip + salt))
	return hex.EncodeToString(hash[:])
}
