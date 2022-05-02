package gognet

import (
	"fmt"
	"net"
)

type Conn struct {
	Connection *net.UDPConn
	Addr       net.UDPAddr

	closed chan bool

	packets chan []byte
}

func (connection *Conn) Close() {
	connection.closed <- true
}

func (connection *Conn) Write(bytes []byte) (n int, err error) {
	return connection.Connection.WriteToUDP(bytes, &connection.Addr)
}

func (connection *Conn) Read() (b []byte, err error) {
	select {
	case <-connection.closed:
		connection.Close()
		return nil, fmt.Errorf("connection closed %x", connection.Addr)
	case packet := <-connection.packets:
		return packet, nil
	}
}
