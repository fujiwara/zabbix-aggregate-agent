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
	DefaultTimeout  = 60
	DefaultAddress  = "127.0.0.1:10052"
	Debug           = 0
	Info            = 1
	Error           = 2
	DefaultLogLevel = Info
)

var LogLabel = []string{"DEBUG", "INFO", "ERROR"}

type Agent struct {
	Name          string
	Listen        string
	ListGenerator func() ([]string, error)
	Timeout       int
	LogPrefix     string
	MinLogLevel   int
}

func NewAgent(name string, listen string, timeout int) (a *Agent) {
	if listen == "" {
		listen = DefaultAddress
	}
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	return &Agent{
		Listen:      listen,
		Timeout:     timeout,
		Name:        name,
		LogPrefix:   "[" + name + "]",
		MinLogLevel: DefaultLogLevel,
	}
}

func (a *Agent) Log(level int, params ...interface{}) {
	if a.MinLogLevel > level {
		return
	}
	var args []interface{}
	args = append(args, LogLabel[level])
	args = append(args, a.LogPrefix)
	args = append(args, params...)
	log.Println(args...)
}

func (a *Agent) Run() (err error) {
	a.Log(Info, "Listing", a.Listen)
	listener, err := net.Listen("tcp", a.Listen)
	if err != nil {
		return
	}
	a.Log(Info, "Ready for connection")
	var conn net.Conn
	for {
		conn, err = listener.Accept()
		if err != nil {
			a.Log(Error, "Error accept:", err)
			continue
		}
		a.Log(Debug, "Accepted connection from", conn.RemoteAddr())
		go a.handleConn(conn)
	}
	return
}

func (a *Agent) handleConn(conn net.Conn) {
	defer a.Log(Debug, "Closing connection:", conn.RemoteAddr())
	defer conn.Close()
	key, err := stream2Data(conn)
	if err != nil {
		a.sendError(conn, err)
		return
	}
	keyString := strings.TrimRight(string(key), "\n")
	a.Log(Debug, "Key:", keyString)
	list, err := a.ListGenerator()
	if err != nil {
		a.sendError(conn, err)
		return
	}
	a.Log(Debug, "List:", list)
	if len(list) == 0 {
		a.sendError(conn, errors.New("Empty list"))
		return
	}
	value, err := a.aggregateValue(list, keyString)
	if err != nil {
		a.sendError(conn, err)
		return
	}
	a.Log(Debug, "Aggregated", keyString, "=", value)
	packet := data2Packet([]byte(value))
	_, err = conn.Write(packet)
	if err != nil {
		a.Log(Error, "Error write:", err)
	}
	return
}

func (a *Agent) sendError(conn net.Conn, err error) {
	a.Log(Error, err)
	packet := data2Packet([]byte(ErrorMessage))
	conn.Write(packet)
}

func (a *Agent) getAsync(host string, key string, timeout int, ch chan []byte) {
	start := time.Now()
	a.Log(Debug, "Sending key:", key, "to", host)
	v, err := Get(host, key, timeout)
	end := time.Now()
	elapsed := int64(end.Sub(start) / time.Millisecond) // msec
	if err != nil {
		a.Log(Error, err)
		v = []byte("0")
	} else {
		a.Log(Debug, "Replied from", host, "in", elapsed, "msec:", string(v))
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
