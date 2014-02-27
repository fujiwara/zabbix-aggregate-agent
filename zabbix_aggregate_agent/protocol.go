package zabbix_aggregate_agent

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

const (
	HeaderString     = "ZBXD"
	HeaderLength     = len(HeaderString)
	HeaderVersion    = uint8(1)
	DataLengthOffset = int64(HeaderLength + 1)
	DataLengthSize   = int64(8)
	DataOffset       = int64(DataLengthOffset + DataLengthSize)
	ErrorMessage     = "ZBX_NOTSUPPORTED"
)

var Terminator = []byte("\n")
var HeaderBytes = []byte(HeaderString)

func data2Packet(data []byte) []byte {
	buf := new(bytes.Buffer)
	buf.Write(HeaderBytes)
	binary.Write(buf, binary.LittleEndian, HeaderVersion)
	binary.Write(buf, binary.LittleEndian, int64(len(data)))
	buf.Write(data)
	return buf.Bytes()
}

func packet2Data(packet []byte) (data []byte, err error) {
	var dataLength int64
	if len(packet) < int(DataOffset) {
		err = errors.New("zabbix protocol packet too short")
		return
	}
	buf := bytes.NewReader(packet[DataLengthOffset:DataOffset])
	err = binary.Read(buf, binary.LittleEndian, &dataLength)
	if err != nil {
		return
	}
	data = packet[DataOffset : DataOffset+dataLength]
	return
}

func stream2Data(conn net.Conn) (rdata []byte, err error) {
	// read header "ZBXD\x01"
	head := make([]byte, DataLengthOffset)
	_, err = conn.Read(head)
	if err != nil {
		return
	}
	if bytes.Equal(head[0:HeaderLength], HeaderBytes) && head[HeaderLength] == byte(HeaderVersion) {
		rdata, err = parseBinary(conn)
	} else {
		rdata, err = parseText(conn, head)
	}
	return
}

func parseBinary(conn net.Conn) (rdata []byte, err error) {
	// read data length
	var dataLength int64
	err = binary.Read(conn, binary.LittleEndian, &dataLength)
	if err != nil {
		return
	}
	// read data body
	buf := make([]byte, 1024)
	data := new(bytes.Buffer)
	total := 0
	size := 0
	for total < int(dataLength) {
		size, err = conn.Read(buf)
		if err != nil {
			return
		}
		if size == 0 {
			break
		}
		total = total + size
		data.Write(buf[0:size])
	}
	rdata = data.Bytes()
	return
}

func parseText(conn net.Conn, head []byte) (rdata []byte, err error) {
	data := new(bytes.Buffer)
	data.Write(head)
	buf := make([]byte, 1024)
	size := 0
	for {
		// read data while "\n" found
		size, err = conn.Read(buf)
		if err != nil {
			return
		}
		if size == 0 {
			break
		}
		i := bytes.Index(buf[0:size], Terminator)
		if i == -1 {
			// terminator not found
			data.Write(buf[0:size])
			continue
		} else {
			// terminator found
			data.Write(buf[0:i])
			break
		}
	}
	rdata = data.Bytes()
	return
}
