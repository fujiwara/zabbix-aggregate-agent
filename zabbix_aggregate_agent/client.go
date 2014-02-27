package zabbix_aggregate_agent

import (
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
	value, err = Stream2Data(conn)
	return
}
