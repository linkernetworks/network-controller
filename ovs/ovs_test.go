package ovs

import (
	"github.com/digitalocean/go-openvswitch/ovs"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
)

var ovsFile = "/usr/bin/ovs-vsctl"

func changeVSCtl(t *testing.T) os.FileMode {
	info, err := os.Stat(ovsFile)
	assert.NoError(t, err)

	err = os.Chmod(ovsFile, 0666)
	assert.NoError(t, err)
	return info.Mode()
}

func resetVSCtl(mode os.FileMode) {
	os.Chmod(ovsFile, mode)
}

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

func TestAddBridgeFail(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	mode := changeVSCtl(t)
	defer resetVSCtl(mode)
	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	assert.Error(t, err)
	c := ovs.New(ovs.Sudo())
	defer c.VSwitch.DeleteBridge(bridgeName)
}

func TestAddDelPort(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	c := ovs.New(ovs.Sudo())
	defer c.VSwitch.DeleteBridge(bridgeName)

	linkName := "test0"
	err = exec.Command("ip", "link", "add", linkName, "type", "veth", "peer", "name", linkName+"_peer").Run()
	assert.NoError(t, err)
	defer exec.Command("ip", "link", "del", linkName).Output()
	err = AddPort(bridgeName, linkName)
	assert.NoError(t, err)

	br, err := c.VSwitch.PortToBridge(linkName)
	assert.NoError(t, err)
	assert.Equal(t, br, bridgeName)

	err = DeletePort(bridgeName, linkName)
	assert.NoError(t, err)

	_, err = c.VSwitch.PortToBridge(linkName)
	assert.Error(t, err)

}

func TestAddDelPortFail(t *testing.T) {
	mode := changeVSCtl(t)
	defer resetVSCtl(mode)

	bridgeName := "bridge0"
	err := AddPort(bridgeName, "0")
	assert.Error(t, err)
	err = DeletePort(bridgeName, "0")
	assert.Error(t, err)
}
