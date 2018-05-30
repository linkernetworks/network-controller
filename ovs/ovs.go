package ovs

import (
	"fmt"

	"github.com/digitalocean/go-openvswitch/ovs"
)

// ovs-vsctl add-port br0 eth0
func AddPort(bridgeName, ifName string) error {
	c := ovs.New(ovs.Sudo())
	if err := c.VSwitch.AddPort(bridgeName, ifName); err != nil {
		return fmt.Errorf("failed to add port: %v", err)
	}
	return nil
}
