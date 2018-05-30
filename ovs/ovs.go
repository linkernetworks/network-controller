package ovs

import (
	"fmt"

	"github.com/digitalocean/go-openvswitch/ovs"
)

// ovs-vsctl --may-exist add-br ovsbr0
func AddBridge(bridgeName string) error {
	c := ovs.New(ovs.Sudo())
	if err := c.VSwitch.AddBridge(bridgeName); err != nil {
		return fmt.Errorf("failed to add bridge: %v", err)
	}
	return nil
}

// ovs-vsctl add-port br0 eth0
func AddPort(bridgeName, ifName string) error {
	c := ovs.New(ovs.Sudo())
	if err := c.VSwitch.AddPort(bridgeName, ifName); err != nil {
		return fmt.Errorf("failed to add port: %v", err)
	}
	return nil
}

// ovs-vsctl del-port br0 eth0
func DeletePort(bridgeName, ifName string) error {
	c := ovs.New(ovs.Sudo())
	if err := c.VSwitch.DeletePort(bridgeName, ifName); err != nil {
		return fmt.Errorf("failed to delete port: %v", err)
	}
	return nil
}
