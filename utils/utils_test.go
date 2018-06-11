package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVethName(t *testing.T) {
	o := GenerateVethName("12345678", "12345678")
	assert.Equal(t, "veth33cdbc38", o)
}

func TestVerifyIP(t *testing.T) {
	o := VerifyIP("8.8.8.8")
	assert.Equal(t, true, o)
}

func TestVerifyCIDR(t *testing.T) {
	o := VerifyCIDR("10.0.0.1/24")
	assert.Equal(t, true, o)
}
