package xirrusswitch

import (
	"github.com/ziutek/telnet"
	"time"
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

func sendln(t *telnet.Conn, s string) error {
	err := t.SetWriteDeadline(time.Now().Add(timeout))
	if err != nil {
		return err
	}
	buf := make([]byte, len(s)+1)
	copy(buf, s)
	buf[len(s)] = '\n'
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
	t.SetUnixWriteMode(true)

	_, err = login(t, username, password)

	return t, err
}
