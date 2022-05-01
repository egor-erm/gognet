package network

import (
	"bytes"
	"encoding/binary"
)

type OpenConnectionReply1 struct {
	ServerGUID int32
}

func (pk *OpenConnectionReply1) Write(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.BigEndian, IDOpenConnectionReply1)
	_ = binary.Write(buf, binary.BigEndian, pk.ServerGUID)
}

func (pk *OpenConnectionReply1) Read(buf *bytes.Buffer) error {
	return binary.Read(buf, binary.BigEndian, &pk.ServerGUID)
}
