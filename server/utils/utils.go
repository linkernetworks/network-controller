package utils

import (
	"crypto/sha256"
	"encoding/hex"
)

func HashUUID(UUID string) string {
	hash := sha256.New()
	hash.Write([]byte(UUID))
	md := hash.Sum(nil)
	mdStr := hex.EncodeToString(md)
	return mdStr
}
