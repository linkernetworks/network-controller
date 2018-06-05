package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateVethName(t *testing.T) {
	o := GenerateVethName("12345678")
	assert.Equal(t, "vethef797c81", o)
}
