package xirrusswitch

import (
	"github.com/ziutek/telnet"
	"time"
	"fmt"
	"bytes"
	"strings"
)

////////////////////////////////////////////////////////////
// Internal Methods

func expect(t *telnet.Conn, d ...string) error {
	err := t.SetReadDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}
	err = t.SkipUntil(d...)
	if err != nil {
		return err
	}
	return nil
}

func expectLong(t *telnet.Conn, d ...string) error {
	err := t.SkipUntil(d...)
	if err != nil {
		return err
	}
	return nil
}

func expectLarge(t *telnet.Conn) ([]byte, error) {
	fmt.Println("----")
	var data bytes.Buffer
	for true {
		dataChunk, readErr := t.ReadUntil("--More--", "#")
		if readErr != nil {
			return nil, readErr
		}
		data.Write(dataChunk)
		if strings.Contains(string(dataChunk), "--More--") {
			sendln(t, " ")
			continue
		}
		if strings.Contains(string(dataChunk), "(vlan)#") {
			break
		}
		/*dataChunk, readErr := t.ReadUntil("--More--", "#")
		fmt.Println(string(dataChunk))
		data.Write(dataChunk)
		sendln(t, " ")
		err := t.SkipUntil("--More--","#")
		if readErr != nil {
			return nil, err
		}
		if strings.Contains(string(dataChunk), "(vlan)#") {
			break
		}
		fmt.Println("Large Output, sending space.")*/
	}
	fmt.Println("----")
	return data.Bytes(), nil
}

func sendln(t *telnet.Conn, s string) error {
  return send(t, s+"\n")
}

func send(t *telnet.Conn, s string) error {
	err := t.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	_, err = t.Write(buf)
	if err != nil {
		return err
	}
	return nil
}
func login(t *telnet.Conn, user string, passwd string) (data []byte, err error) {
	expect(t, "name: ")
	sendln(t, user)
	expect(t, "ssword: ")
	sendln(t, passwd)
	expect(t, "#")
	sendln(t, "")
	data, err = t.ReadBytes('#')
	return data, err
}

func GetTelnetToSwitch(address string, username string, password string) (t *telnet.Conn, err error) {
	t, err = telnet.Dial("tcp", address+":23")
	if err != nil {
		return nil, err
	}
	t.SetUnixWriteMode(true)

	_, err = login(t, username, password)
	if err != nil {
		return nil, err
	}

	return t, err
}
