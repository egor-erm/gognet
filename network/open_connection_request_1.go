package network

import (
	"bytes"
)

type OpenConnectionRequest1 struct {
	protocol byte
}

func (pk *OpenConnectionRequest1) Write(buf *bytes.Buffer) {
	buf.Write([]byte{Protocol_Version})
}

func (pk *OpenConnectionRequest1) Read(buf *bytes.Buffer) error {
	pk.protocol, _ = buf.ReadByte()
	_, err := buf.ReadByte()
	return err
}
