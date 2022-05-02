package gognet

import (
	"fmt"
	"net"
)

type Conn struct {
	Connection *net.UDPConn
	Addr       net.UDPAddr

	closed chan *Conn

	packets chan []byte
}

func (connection *Conn) Close() {
	connection.closed <- connection
}

func (connection *Conn) Write(bytes []byte) (n int, err error) {
	return connection.Connection.WriteToUDP(bytes, &connection.Addr)
}

func (connection *Conn) Read() (b []byte, err error) {
	select {
	case packet := <-connection.packets:
		return packet, nil
	case <-connection.closed:
		return nil, fmt.Errorf("connection closed %x", connection.Addr)
	}
}
