package openvswitch

import (
	"fmt"

	"github.com/linkernetworks/go-openvswitch/ovs"
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

// CreateBridge is a function for create bridge
// userspace datapath
// ovs-vsctl add-br br0 -- set bridge br0 datapath_type=netdev
// kernel datapath
// ovs-vsctl add-br br1 -- set bridge br1 datapath_type=system
func (o *OVSManager) CreateBridge(bridgeName, dpType string) error {
	if err := o.Client.VSwitch.AddBridgeWithType(bridgeName, dpType); err != nil {
		return fmt.Errorf("Failed to add bridge %s. Datapath type %s: %v", bridgeName, dpType, err)
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
	return bridges, nil
}

// AddPort : ovs-vsctl add-port br0 eth0
func (o *OVSManager) AddPort(bridgeName, ifName string) error {
	if err := o.Client.VSwitch.AddPort(bridgeName, ifName); err != nil {
		return fmt.Errorf("Failed to add port: %s on %s: %v", ifName, bridgeName, err)
	}
	return nil
}

// AddDPDKPort : ovs-vsctl add-port br0 dpdk0 -- set Interface dpdk0 type=dpdk options:dpdk-devargs=0000:00:08.0
func (o *OVSManager) AddDPDKPort(bridgeName, ifName, dpdkDevargs string) error {
	if err := o.Client.VSwitch.AddDPDKPort(bridgeName, ifName, dpdkDevargs); err != nil {
		return fmt.Errorf("Failed to add dpdk port: %s on %s dpdkDevargs: %s: %v", ifName, bridgeName, dpdkDevargs, err)
	}
	return nil
}

// GetPort : ovs-vsctl --format=json get port eth0 tag vlan_mode trunk
func (o *OVSManager) GetPort(ifName string) (ovs.PortOptions, error) {
	portOptions, err := o.Client.VSwitch.Get.Port(ifName)
	if err != nil {
		return ovs.PortOptions{}, fmt.Errorf("Failed to get port options: %v", err)
	}
	return portOptions, nil
}

// SetPort : ovs-vsctl --format=json set port eth0 vlan_mode=trunk trunk=1,2,3,4,5
func (o *OVSManager) SetPort(ifName string, portOptions ovs.PortOptions) error {
	if err := o.Client.VSwitch.Set.Port(ifName, portOptions); err != nil {
		return fmt.Errorf("Failed to set port options: %v on %s: %v", portOptions, ifName, err)
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

// DumpPorts : ovs-ofctl dump-ports br0
func (o *OVSManager) DumpPorts(bridgeName string) ([]*ovs.PortStats, error) {
	ports, err := o.Client.OpenFlow.DumpPorts(bridgeName)
	if err != nil {
		return ports, fmt.Errorf("Failed to dump ports: %v", err)
	}
	return ports, nil
}

// DumpPort : ovs-ofctl dump-ports br0 eth0
func (o *OVSManager) DumpPort(bridgeName, portName string) (*ovs.PortStats, error) {
	port, err := o.Client.OpenFlow.DumpPort(bridgeName, portName)
	if err != nil {
		return port, fmt.Errorf("Failed to dump port: %v", err)
	}
	return port, nil
}
