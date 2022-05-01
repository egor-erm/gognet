package network

import "bytes"

const (
	Protocol_Version byte = 0x01

	IDOpenConnectionRequest1      byte = 0x01
	IDOpenConnectionReply1        byte = 0x02
	IDDisconnectNotification      byte = 0x03
	IDIncompatibleProtocolVersion byte = 0x04
)

type Packet interface {
	Read(*bytes.Buffer)
	Write(*bytes.Buffer) error
}