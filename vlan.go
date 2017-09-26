package xirrusswitch

import (
	"github.com/ziutek/telnet"
	"strconv"
	"strings"
	"fmt"
)

/*
 * Raw VLAN Support
 */
func processVLANList(data []byte) []VLAN {
	dataString := string(data[:])                  //Convert bytes to string
	dataLines := strings.Split(dataString, "\r\n") //Convert to string array

	var vlans []VLAN
	for _, line := range dataLines {
		fields := strings.Fields(line)                     // Split line into array delimited by whitespace
		if _, err := strconv.Atoi(fields[0]); err == nil { //Check if the line starts with a number, if it does, its a VLAN
			vlan := VLAN{}
			vlan.Id, _ = strconv.Atoi(fields[0])
			if len(fields) == 4 {
				vlan.Name = fields[1]
				vlan.Assignment = fields[3]
			} else {
				vlan.Name = ""
				vlan.Assignment = fields[2]
			}
			vlans = append(vlans, vlan)
		}
	}
	return vlans
}

//GetVLANS Get VLANs
func GetVLANS(t *telnet.Conn) ([]VLAN, error) {
	sendln(t, "vlan")
	sendln(t, "sh vlan")
	expectLong(t, "#")
	data, err := t.ReadBytes('#')
	if err != nil {
		return nil, err
	}
	vlans := processVLANList(data)
	sendln(t, "exit")
	expect(t, "#")
	return vlans, err

}

//CreateVLAN Creates a VLAN
func CreateVLAN(t *telnet.Conn, vlan VLAN) {
	//Chop off extra bits
	if vlan.Assignment[0] == ',' {
		vlan.Assignment = vlan.Assignment[1:]
	}
	if vlan.Assignment[len(vlan.Assignment)-1] == ',' {
		vlan.Assignment = vlan.Assignment[:len(vlan.Assignment)-1]
	}
	sendln(t, "vlan")
	expect(t, "#")
	sendln(t, fmt.Sprintf("tag-group %d %s %s", vlan.Id, vlan.Name, vlan.Assignment))
	sendln(t,"exit")
	expect(t, "#")
}

//GetVLAN Get a VLAN by ID
func GetVLAN(t *telnet.Conn, vlanid int) (*VLAN, error){
	vlans, err := GetVLANS(t)
	if err != nil {
		return nil, err
	}

	//Find one that matches
	var existingVLAN VLAN
	for i, vlan := range vlans {
		if vlan.Id == vlanid {
			existingVLAN = vlans[i]
		}
	}

	return &existingVLAN, err
}

/*
 * VLANPortConfig Logic
 */

func processVLANPortConfigList(data []byte) []VLANPortConfig {
	dataString := string(data[:])                  //Convert bytes to string
	dataLines := strings.Split(dataString, "\r\n") //Convert to string array

	var portConfigs []VLANPortConfig
	for _, line := range dataLines {
		fields := strings.Fields(line)                     // Split line into array delimited by whitespace
		//Ignore bad lines
		if len(fields) != 6 {
			continue
		}
		if _, err := strconv.Atoi(fields[0]); err == nil { //Check if the line starts with a number, if it does, its a VLANPortConfig
			//Start building portconfig struct
			vlanPortConfig := VLANPortConfig{}
			vlanPortConfig.Port, _ = strconv.Atoi(fields[0])
			vlanPortConfig.NativeVLAN, _ = strconv.Atoi(fields[1])
			vlanPortConfig.FrameType = fields[2]
			vlanPortConfig.IngressFilter = fields[3]
			vlanPortConfig.EgressRule = fields[4]
			vlanPortConfig.PortType = fields[5]

			portConfigs = append(portConfigs, vlanPortConfig)
		}
	}
	return portConfigs
}

//GetVLANPortConfigs Get VLANPortConfigs
func GetVLANPortConfigs(t *telnet.Conn) ([]VLANPortConfig, error) {
	sendln(t, "vlan")
	expect(t, "#")
	sendln(t, "sh port-config")
	data, err := expectLarge(t)
	if err != nil {
		return nil, err
	}
	vlans := processVLANPortConfigList(data)
	sendln(t, "exit")
	expect(t, "#")
	return vlans, err
}