package util

import (
	"crypto/sha256"
	"encoding/hex"
)

func Sha256String(str string) string {
	hashedByte := sha256.Sum256([]byte(str))
	hashedString := hex.EncodeToString(hashedByte[:])
	return hashedString
}
