package def

import (
	"net"
)

/*
	use when node register to master
*/
type RegisterToMasterResponse struct {
	NodeName string
	NodeID   int
	CniIP    net.IP
}

type RegisterToMasterRequest struct {
	NodeName  string
	LocalIP   net.IP
	LocalPort int
	ProxyPort int
}

/*
	node rend heartbeat to master, no response
*/
type NodeToMasterHeartBeatRequest struct {
	NodeID       int
	PodInstances []PodInstance
}
