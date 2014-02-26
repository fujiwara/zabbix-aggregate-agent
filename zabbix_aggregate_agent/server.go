package zabbix_aggregate_agent

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultTimeout = 5
	DefaultAddress = "127.0.0.1:10052"
)

type Agent struct {
	Name          string
	Listen        string
	ListGenerator func(string) []string
	ListSource    string
	Timeout       int
	LogPrefix     string
}

func NewAgent(name string, listen string, timeout int) (a *Agent) {
	if listen == "" {
		listen = DefaultAddress
	}
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	return &Agent{
		Listen:    listen,
		Timeout:   timeout,
		Name:      name,
		LogPrefix: "[" + name + "]",
	}
}

func (a *Agent) Log(params ...interface{}) {
	var args []interface{}
	args = append(args, a.LogPrefix)
	args = append(args, params...)
	log.Println(args...)
}

func (a *Agent) Run() (err error) {
	a.Log("Listing", a.Listen)
	listener, err := net.Listen("tcp", a.Listen)
	if err != nil {
		return
	}
	a.Log("Ready for connection")
	var conn net.Conn
	for {
		conn, err = listener.Accept()
		if err != nil {
			a.Log("Error accept:", err)
			continue
		}
		a.Log("Accepted connection from", conn.RemoteAddr())
		go a.handleConn(conn)
	}
	return
}

func (a *Agent) handleConn(conn net.Conn) {
	defer a.Log("Closing connection:", conn.RemoteAddr())
	defer conn.Close()
	key, err := stream2Data(conn)
	if err != nil {
		a.sendError(conn, err)
		return
	}
	keyString := strings.Trim(string(key), "\n")
	list := a.ListGenerator(a.ListSource)
	value, err := a.aggregateValue(list, keyString)
	if err != nil {
		a.sendError(conn, err)
		return
	}
	a.Log("Aggregated", keyString, ":", value)
	packet := data2Packet([]byte(value))
	_, err = conn.Write(packet)
	if err != nil {
		a.Log("Error write:", err)
	}
	return
}

func (a *Agent) sendError(conn net.Conn, err error) {
	a.Log(conn.RemoteAddr(), err)
	packet := data2Packet([]byte(ErrorMessage))
	conn.Write(packet)
}

func (a *Agent) getAsync(host string, key string, timeout int, ch chan []byte) {
	start := time.Now()
	a.Log("Sending key:", key, "to", host)
	v, err := Get(host, key, timeout)
	end := time.Now()
	elapsed := int64(end.Sub(start) / time.Millisecond) // msec
	if err != nil {
		a.Log(err)
		v = []byte("0")
	} else {
		a.Log("Replied from", host, "in", elapsed, "msec:", string(v))
	}
	ch <- v
}

func (a *Agent) aggregateValue(list []string, key string) (rvalue string, err error) {
	ch := make(chan []byte)
	for _, host := range list {
		go a.getAsync(host, key, a.Timeout, ch)
	}
	errs := 0
	isInt := true
	var value float64
	var valueBuf bytes.Buffer
	for _ = range list {
		v := <-ch
		vs := string(v)
		vf, err := strconv.ParseFloat(string(v), 64)
		if err != nil {
			// may be string
			if vs == ErrorMessage {
				errs++
			} else {
				valueBuf.WriteString(vs)
				valueBuf.WriteString("\n")
			}
			continue
		}
		if strings.Index(vs, ".") != -1 {
			// may be float
			isInt = false
		}
		value = value + vf
	}
	if errs == len(list) {
		err = errors.New("All replied values could not be parsed as float")
		return
	}
	if valueBuf.Len() > 0 {
		return valueBuf.String(), err
	}
	if isInt {
		return fmt.Sprintf("%d", int64(value)), err
	}
	return fmt.Sprintf("%.6f", value), err
}
