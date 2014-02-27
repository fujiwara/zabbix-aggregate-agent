package zabbix_aggregate_agent

import (
	"bytes"
	"net"
	"time"
)

func Get(host string, key string, timeout int) (value []byte, err error) {
	conn, err := net.DialTimeout("tcp", host, time.Duration(timeout)*time.Second)
	if err != nil {
		return
	}
	defer conn.Close()
	msg := Data2Packet([]byte(key))
	_, err = conn.Write(msg)
	if err != nil {
		return
	}
	reply := make([]byte, 1024)
	data := new(bytes.Buffer)
	var size int
	for {
		size, err = conn.Read(reply)
		if size == 0 {
			break
		}
		if err != nil {
			return
		}
		data.Write(reply[0:size])
	}
	value, err = Packet2Data(data.Bytes())
	return
}
