package xirrusswitch

import (
	"github.com/ziutek/telnet"
	"strconv"
	"strings"
)

func processVLANList(data []byte) []VLAN {
	dataString := string(data[:])                  //Convert bytes to string
	dataLines := strings.Split(dataString, "\r\n") //Convert to string array

	var vlans []VLAN
	for _, line := range dataLines {
		fields := strings.Fields(line)                     // Split line into array delimited by whitespace
		if _, err := strconv.Atoi(fields[0]); err == nil { //Check if the line starts with a number, if it does, its a VLAN
			vlan := VLAN{}
			vlan.id, _ = strconv.Atoi(fields[0])
			if len(fields) == 4 {
				vlan.name = fields[1]
				vlan.assignment = fields[3]
			} else {
				vlan.name = ""
				vlan.assignment = fields[2]
			}
			vlans = append(vlans, vlan)
		}
	}
	return vlans
}

func GetVLANS(t *telnet.Conn) ([]VLAN, error) {
	sendln(t, "vlan")
	sendln(t, "sh vlan")
	expectLong(t, "#")
	data, err := t.ReadBytes('#')
	if err != nil {
		return nil, err
	}
	vlans := processVLANList(data)

	return vlans, err

}
