package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

// Hash 哈希加密
func Hash(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}
