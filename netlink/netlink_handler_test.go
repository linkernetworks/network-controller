package netlinkEvent

import (
	"os"
	"os/exec"
	"testing"

	ovs "github.com/linkernetworks/network-controller/openvswitch"
	"github.com/stretchr/testify/assert"
	"github.com/vishvananda/netlink"
)

var o *ovs.OVSManager

func TestMain(m *testing.M) {
	if _, ok := os.LookupEnv("TEST_OVS"); ok {
		o = ovs.New()
		retCode := m.Run()
		os.Exit(retCode)
	}
}

func TestBridgeOperations(t *testing.T) {
	bridges, err := o.ListBridges()
	assert.NoError(t, err)
	assert.Equal(t, 0, len(bridges))

	bridgeName := "br0"
	err = o.CreateBridge(bridgeName)
	defer o.DeleteBridge(bridgeName)
	assert.NoError(t, err)

	hName := "test0"
	cName := "test0_peer"
	err = exec.Command("ip", "link", "add", hName, "type", "veth", "peer", "name", cName).Run()
	assert.NoError(t, err)
	defer exec.Command("ip", "link", "del", hName).Output()
	err = o.AddPort(bridgeName, hName)
	assert.NoError(t, err)

	ports, err := o.ListPorts(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(ports))

	link := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			MTU:         1400,
			Name:        hName,
			MasterIndex: 1,
			ParentIndex: 0,
		},
	}
	nl := netlink.LinkUpdate{
		Link: link,
	}

	//case 1. not a orphn veth (don't do anthing)
	ret := RemoveVethFromOVS(nl)
	assert.Equal(t, false, ret)

	nl.Link.Attrs().MasterIndex = 0
	//Case 2, masterIndex and PartentIndex is 0, need to remove from ovs
	ret = RemoveVethFromOVS(nl)
	assert.Equal(t, false, ret)

	ports, err = o.ListPorts(bridgeName)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(ports))
}
