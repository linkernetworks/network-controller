package openvswitch

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

// ovs-vsctl del-br ovsbr0
func DeleteBridge(bridgeName string) error {
	c := ovs.New(ovs.Sudo())
	if err := c.VSwitch.DeleteBridge(bridgeName); err != nil {
		return fmt.Errorf("failed to delete bridge: %v", err)
	}
	return nil
}

// ovs-vsctl list-br
func ListBridges() ([]string, error) {
	c := ovs.New(ovs.Sudo())
	bridges, err := c.VSwitch.ListBridges()
	if err != nil {
		return bridges, fmt.Errorf("failed to list bridges: %v", err)
	}
	return bridges, nil
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

// ovs-vsctl list-ports br0
func ListPorts(bridgeName string) ([]string, error) {
	c := ovs.New(ovs.Sudo())
	ports, err := c.VSwitch.ListPorts(bridgeName)
	if err != nil {
		return ports, fmt.Errorf("failed to list ports: %v", err)
	}
	return ports, nil
}

// ovs-ofctl add-flow br0 "flow"
func AddFlow(bridgeName string, flowString string) error {
	flow := &ovs.Flow{}
	flow.UnmarshalText([]byte(flowString))
	c := ovs.New(ovs.Sudo())
	if err := c.OpenFlow.AddFlow(bridgeName, flow); err != nil {
		return fmt.Errorf("failed to add flow: %v", err)
	}
	return nil
}

// ovs-ofctl del-flow br0 "flow"
func DeleteFlow(bridgeName string, flowString string) error {
	flow := &ovs.Flow{}
	flow.UnmarshalText([]byte(flowString))
	c := ovs.New(ovs.Sudo())
	if err := c.OpenFlow.DelFlows(bridgeName, flow.MatchFlow()); err != nil {
		return fmt.Errorf("failed to delete flows: %v", err)
	}
	return nil
}

// ovs-ofctl dump-flows br0
func DumpFlows(bridgeName string) ([]*ovs.Flow, error) {
	c := ovs.New(ovs.Sudo())
	flows, err := c.OpenFlow.DumpFlows(bridgeName)
	if err != nil {
		return flows, fmt.Errorf("failed to dump flows: %v", err)
	}
	return flows, nil
}
