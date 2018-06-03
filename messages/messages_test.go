package messages

import (
	"github.com/golang/protobuf/jsonpb"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshalAddPortRequestJson(t *testing.T) {
	j := `
	{
		"bridgeName":"br0",
		"ifaceName" :"eth"
	}
	`
	req := AddPortRequest{}
	err := jsonpb.UnmarshalString(j, &req)
	assert.NoError(t, err)
}

func TestUnmarshalDeletePortRequestJson(t *testing.T) {
	j := `
	{
		"bridgeName":"br0",
		"ifaceName" :"eth"
	}
	`
	req := DeletePortRequest{}
	err := jsonpb.UnmarshalString(j, &req)
	assert.NoError(t, err)
}

func TestUnmarshalAddFlowRequestJson(t *testing.T) {
	j := `
	{
		"bridgeName":"br0",
		"flowString" :"ip,tcp,action=drop"
	}
	`
	req := AddFlowRequest{}
	err := jsonpb.UnmarshalString(j, &req)
	assert.NoError(t, err)
}

func TestUnmarshalDeleteFlowRequestJson(t *testing.T) {
	j := `
	{
		"bridgeName":"br0",
		"flowString" :"ip,tcp,action=drop"
	}
	`
	req := DeleteFlowRequest{}
	err := jsonpb.UnmarshalString(j, &req)
	assert.NoError(t, err)
}
