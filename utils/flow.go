package utils

import (
  "github.com/digitalocean/go-openvswitch/ovs"

  "strings"
)

func ConvertOVSFlow(f map[string]interface{}) *ovs.Flow {

  flow := &ovs.Flow{}
  for key, value := range f {

    switch key {
    case "priority":
      flow.Priority = value.(int)

    case "protocol":
      flow.Protocol = ovs.Protocol(value.(string))

    case "in_port":
      flow.InPort = value.(int)

    case "table":
      flow.Table = value.(int)

    case "idle_timeout":
      flow.IdleTimeout = value.(int)

    case "cookie":
      flow.Cookie = uint64(f["cookie"].(float64))

    case "actions":
      actions := make([]ovs.Action, len(value.([]interface{})))
      for i, action := range value.([]interface{}) {
        switch strings.ToLower(action.(string)) {
        case "drop":
          actions[i] = ovs.Drop()
        case "flood":
          actions[i] = ovs.Flood()
        case "in_port":
          actions[i] = ovs.InPort()
        case "local":
          actions[i] = ovs.Local()
        case "normal":
          actions[i] = ovs.Normal()
        case "strip_vlan":
          actions[i] = ovs.StripVLAN()
        }
      }
      flow.Actions = actions
    }
  }

  return flow
}
