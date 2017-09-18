package xirrusswitch

import "time"

const timeout = 10 * time.Second

////////////////////////////////////////////////////////////
// Structs

//VLAN Struct that wraps around Xirrus VLANs
type VLAN struct {
	Id         int
	Name       string
	Assignment string
}

//LLDPPeer Struct that wraps around Xirrus LLDP neighbors
type LLDPPeer struct {
	localPort          string
	chassisID          string
	portID             string
	portDescription    string
	systemName         string
	systemDescription  string
	systemCapabilities []string
	managementAddress  string
}
