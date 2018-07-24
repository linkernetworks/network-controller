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

// GenerateVethName : Use pod-uuid and container-vethName to generate vethXXXXXXXX for bridge's interface name
func GenerateVethName(podUUID, containerVethName string) string {
	str := sha256String(podUUID + containerVethName)
	return fmt.Sprintf("veth%s", str[0:8])
}

// IsValidVLANTag : Check VLAN tag is valided
func IsValidVLANTag(tag int32) bool {
	if tag < 0 || tag > 4095 {
		return false
	}
	return true
}
