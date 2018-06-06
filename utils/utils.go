package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func sha256String(str string) string {
	hash := sha256.New()
	hash.Write([]byte(str))
	md := hash.Sum(nil)
	return fmt.Sprintf("%s", hex.EncodeToString(md))
}

// GenerateVethName : Use pod-uuid to generate vethXXXXXXXX for bridge's interface name
func GenerateVethName(uuid string) string {
	str := sha256String(uuid)
	return fmt.Sprintf("veth%s", str[0:8])
}
