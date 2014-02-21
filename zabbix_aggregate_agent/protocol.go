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

func data2Packet(data []byte) []byte {
	buf := new(bytes.Buffer)
	buf.WriteString(HeaderString)
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
	headBuf := make([]byte, DataLengthOffset)
	_, err = conn.Read(headBuf)
	if err != nil {
		return
	}
	if string(headBuf[0:HeaderLength]) != HeaderString || headBuf[HeaderLength] != byte(HeaderVersion) {
		err = errors.New("invalid header")
		return
	}

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
