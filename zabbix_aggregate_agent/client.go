package zabbix_aggregate_agent

import (
	"net"
	"strings"
	"time"
)

func IsIncludesPortNumber (addr string) bool {
	if strings.Index(addr, "[") == 0 && strings.Index(addr, "]") == len(addr) - 1 {
		// ipv6
		return false
	}
	if strings.Index(addr, ":") == -1 {
		return false
	}
	return true
}

func Get(host string, key string, timeout int) (value []byte, err error) {
	if ! IsIncludesPortNumber(host) {
		host = host + ":10050"
	}
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
