package gognet

import (
	"fmt"
	"net"
)

type Conn struct {
	conn *net.UDPConn
	addr net.UDPAddr

	closed chan struct{}

	packets chan []byte
}

func (connection *Conn) Write(bytes []byte) (n int, err error) {
	return connection.conn.WriteToUDP(bytes, &connection.addr)
}

func (connection *Conn) Read() (b []byte, err error) {
	select {
	case packet := <-connection.packets:
		return packet, nil
	case <-connection.closed:
		return nil, fmt.Errorf("connection closed %x", connection.addr)
	}
}
