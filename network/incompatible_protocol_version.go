package network

import (
	"bytes"
	"encoding/binary"
)

type IncompatibleProtocolVersion struct {
	ServerProtocol byte
	ServerGUID     int32
}

func (pk *IncompatibleProtocolVersion) Write(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.BigEndian, IDIncompatibleProtocolVersion)
	_ = binary.Write(buf, binary.BigEndian, pk.ServerProtocol)
	_ = binary.Write(buf, binary.BigEndian, pk.ServerGUID)
}

func (pk *IncompatibleProtocolVersion) Read(buf *bytes.Buffer) error {
	_ = binary.Read(buf, binary.BigEndian, &pk.ServerProtocol)
	return binary.Read(buf, binary.BigEndian, &pk.ServerGUID)
}
