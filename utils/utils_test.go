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

func TestIsInvalidIP(t *testing.T) {
	o := IsValidIP("abcd.efgh.ijkl.mnop")
	assert.Equal(t, false, o)
}

func TestIsValidCIDR(t *testing.T) {
	o := IsValidCIDR("10.0.0.1/24")
	assert.Equal(t, true, o)
}

func TestIsInvalidCIDR(t *testing.T) {
	o := IsValidCIDR("0.0.0.0")
	assert.Equal(t, false, o)
}
