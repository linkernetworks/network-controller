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
