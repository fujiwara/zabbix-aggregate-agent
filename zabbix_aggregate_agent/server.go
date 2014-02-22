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
)

type Agent struct {
	ListGenerator func(string) []string
	ListSource    string
	Timeout       int
}

func NewAgent() (a *Agent) {
	return &Agent{Timeout: DefaultTimeout}
}

func (a *Agent) Run(bindAddress string) (err error) {
	log.Println("Listing", bindAddress)
	listener, err := net.Listen("tcp", bindAddress)
	if err != nil {
		return
	}
	log.Println("Ready for connection")
	var conn net.Conn
	for {
		conn, err = listener.Accept()
		if err != nil {
			log.Println("Error accept:", err)
			continue
		}
		log.Println("Accepted connection from", conn.RemoteAddr())
		go a.handleConn(conn)
	}
	return
}

func (a *Agent) handleConn(conn net.Conn) {
	defer conn.Close()
	key, err := stream2Data(conn)
	if err != nil {
		sendError(conn, err)
		return
	}
	list := a.ListGenerator(a.ListSource)
	value, err := aggregateValue(list, string(key), a.Timeout)
	if err != nil {
		sendError(conn, err)
		return
	}
	packet := data2Packet([]byte(value))
	_, err = conn.Write(packet)
	if err != nil {
		log.Println("Error write:", err)
	}
	return
}

func sendError(conn net.Conn, err error) {
	log.Println(err)
	packet := data2Packet([]byte(ErrorMessage))
	conn.Write(packet)
}

func getAsync(host string, key string, timeout int, ch chan []byte) {
	start := time.Now()
	log.Printf("getting from %s key: %s", host, key)
	v, err := Get(host, key, timeout)
	end := time.Now()
	elapsed := int64(end.Sub(start) / time.Millisecond) // msec
	if err != nil {
		log.Println(err)
		v = []byte("0")
	} else {
		log.Printf("replied from %s in %d msec: %s", host, elapsed, v)
	}
	ch <- v
}

func aggregateValue(list []string, key string, timeout int) (rvalue string, err error) {
	ch := make(chan []byte)
	for _, host := range list {
		go getAsync(host, key, timeout, ch)
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
