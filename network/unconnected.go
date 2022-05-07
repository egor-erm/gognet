package network

import (
	"bytes"
	"encoding/binary"
)

type Unconnected struct {
}

func (pk *Unconnected) Write(buf *bytes.Buffer) {
	_ = binary.Write(buf, binary.BigEndian, IDUnconnected)
}
