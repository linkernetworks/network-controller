package openvswitch

import (
	"fmt"

	"github.com/digitalocean/go-openvswitch/ovs"
)

// OVSManager : contains the client for control ovs-vsctl
type OVSManager struct {
	Client *ovs.Client
}

// New : init OVSManager and the client to be super user
func New() *OVSManager {
	return &OVSManager{
		Client: ovs.New(ovs.Sudo()),
	}
}

// CreateBridge : ovs-vsctl add-br br0
func (o *OVSManager) CreateBridge(bridgeName string) error {
	if err := o.Client.VSwitch.AddBridge(bridgeName); err != nil {
		return fmt.Errorf("Failed to add bridge %s: %v", bridgeName, err)
	}

	return nil
}

// DeleteBridge : ovs-vsctl del-br br0
func (o *OVSManager) DeleteBridge(bridgeName string) error {
	if err := o.Client.VSwitch.DeleteBridge(bridgeName); err != nil {
		return fmt.Errorf("Failed to delete bridge %s: %v", bridgeName, err)
	}
	return nil
}

// ListBridges : ovs-vsctl list-br
func (o *OVSManager) ListBridges() ([]string, error) {
	bridges, err := o.Client.VSwitch.ListBridges()
	if err != nil {
		return bridges, fmt.Errorf("Failed to list bridges: %v", err)
	}

	//FIXME remove this after the upstream has fix the problem
	if len(bridges) == 1 && bridges[0] == "" {
		return []string{}, nil
	}
	return bridges, nil
}

// AddPort : ovs-vsctl add-port br0 eth0
func (o *OVSManager) AddPort(bridgeName, ifName string) error {
	if err := o.Client.VSwitch.AddPort(bridgeName, ifName); err != nil {
		return fmt.Errorf("Failed to add port: %s on %s: %v", ifName, bridgeName, err)
	}
	return nil
}

// DeletePort : ovs-vsctl del-port br0 eth0
func (o *OVSManager) DeletePort(bridgeName, ifName string) error {
	if err := o.Client.VSwitch.DeletePort(bridgeName, ifName); err != nil {
		return fmt.Errorf("Failed to delete port: %s on %s: %v", ifName, bridgeName, err)
	}
	return nil
}

// ListPorts : ovs-vsctl list-ports
func (o *OVSManager) ListPorts(bridgeName string) ([]string, error) {
	ports, err := o.Client.VSwitch.ListPorts(bridgeName)
	if err != nil {
		return ports, fmt.Errorf("Failed to list ports of bridge %s: %v", bridgeName, err)
	}

	//FIXME remove this after the upstream has fix the problem
	if len(ports) == 1 && ports[0] == "" {
		return []string{}, nil
	}
	return ports, nil
}

// AddFlow : ovs-ofctl add-flow br0 "flow"
func (o *OVSManager) AddFlow(bridgeName, flow string) error {
	f := &ovs.Flow{}
	f.UnmarshalText([]byte(flow))
	if err := o.Client.OpenFlow.AddFlow(bridgeName, f); err != nil {
		return fmt.Errorf("Failed to add flow:%s into %s, %v", flow, bridgeName, err)
	}
	return nil
}

// DeleteFlow : ovs-ofctl del-flow br0 "flow"
func (o *OVSManager) DeleteFlow(bridgeName, flow string) error {
	f := &ovs.Flow{}
	f.UnmarshalText([]byte(flow))
	if err := o.Client.OpenFlow.DelFlows(bridgeName, f.MatchFlow()); err != nil {
		return fmt.Errorf("Failed to delete flows: %s on %s :%v", flow, bridgeName, err)
	}
	return nil
}

// DumpFlows : ovs-ofctl dump-flows br0
func (o *OVSManager) DumpFlows(bridgeName string) ([]*ovs.Flow, error) {
	flows, err := o.Client.OpenFlow.DumpFlows(bridgeName)
	if err != nil {
		return flows, fmt.Errorf("Failed to dump flows: %v", err)
	}
	return flows, nil
}
