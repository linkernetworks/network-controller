package ovs

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
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
	defer DeleteBridge(bridgeName)
	assert.NoError(t, err)

	bridges, err := ListBridges()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bridges))
}

func TestDeleteBridge(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	assert.NoError(t, err)

	bridges, err := ListBridges()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(bridges))

	err = DeleteBridge(bridgeName)
	assert.NoError(t, err)
}

func TestListBridges(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	_, err := ListBridges()
	assert.NoError(t, err)
}

func TestAddFlow(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	defer DeleteBridge(bridgeName)
	assert.NoError(t, err)

	flowString := "cookie=1, actions=NORMAL"
	err = AddFlow(bridgeName, flowString)
	assert.NoError(t, err)

	flows, err := DumpFlows(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(flows))
}

func TestDeleteFlows(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	defer DeleteBridge(bridgeName)
	assert.NoError(t, err)

	flowString := "cookie=1, actions=NORMAL"
	err = AddFlow(bridgeName, flowString)
	assert.NoError(t, err)

	err = DeleteFlows(bridgeName, flowString)
	assert.NoError(t, err)

	flows, err := DumpFlows(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(flows))
}

func TestDumpFlows(t *testing.T) {
	if _, ok := os.LookupEnv("TEST_OVS"); !ok {
		t.SkipNow()
	}

	bridgeName := "bridge0"
	err := AddBridge(bridgeName)
	defer DeleteBridge(bridgeName)
	assert.NoError(t, err)

	flows, err := DumpFlows(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(flows))
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
