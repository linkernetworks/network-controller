package ovs

import (
	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/stretchr/testify/assert"
  "github.com/linkernetworks/network-controller/utils"

	"os"
	"testing"
	"encoding/json"
)

func TestAddBridge(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	assert.NoError(t, err)
	c := ovs.New(ovs.Sudo())
	defer c.VSwitch.DeleteBridge(bridgeName)
}

func TestAddFlow(t *testing.T) {
  if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

  bridgeName := "ovs-eth0"
  flowString := `
    {
      "cookie": 1,
      "actions": ["normal"]
    }`
  var flow map[string]interface{}
  err :=  json.Unmarshal([]byte(flowString), &flow)
	assert.NoError(t, err)
  err = AddFlow(bridgeName, utils.ConvertOVSFlow(flow))
	assert.NoError(t, err)
}
