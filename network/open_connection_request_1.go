package network

import (
	"bytes"
	"encoding/binary"
)

type OpenConnectionRequest1 struct {
	Protocol byte
}

func (pk *OpenConnectionRequest1) Write(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.BigEndian, IDOpenConnectionRequest1)
	_ = binary.Write(buf, binary.BigEndian, pk.Protocol)
}

func (pk *OpenConnectionRequest1) Read(buf *bytes.Buffer) error {
	var err error
	pk.Protocol, err = buf.ReadByte()
	return err
}
