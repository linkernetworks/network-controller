package ovs

import (
	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
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

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	assert.NoError(t, err)

	flowString := "cookie=1, actions=NORMAL"
	err = AddFlow(bridgeName, flowString)
	assert.NoError(t, err)

	c := ovs.New(ovs.Sudo())
	flows, err := c.OpenFlow.DumpFlows(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(flows))

	defer c.VSwitch.DeleteBridge(bridgeName)
}

func TestDeleteFlows(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	assert.NoError(t, err)

	flowString := "cookie=1, actions=NORMAL"
	err = AddFlow(bridgeName, flowString)
	assert.NoError(t, err)

	err = DeleteFlows(bridgeName, flowString)
	assert.NoError(t, err)

	c := ovs.New(ovs.Sudo())
	flows, err := c.OpenFlow.DumpFlows(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(flows))

	defer c.VSwitch.DeleteBridge(bridgeName)
}
