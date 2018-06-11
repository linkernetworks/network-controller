package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVethName(t *testing.T) {
	o := GenerateVethName("12345678", "12345678")
	assert.Equal(t, "veth33cdbc38", o)
}

func TestIsValidIP(t *testing.T) {
	o := IsValidIP("8.8.8.8")
	assert.Equal(t, true, o)
}

func TestIsValidCIDR(t *testing.T) {
	o := IsValidCIDR("10.0.0.1/24")
	assert.Equal(t, true, o)
}
