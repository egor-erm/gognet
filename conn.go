package gognet

import (
	"bytes"
	"fmt"
	"net"

	"github.com/egor-erm/gognet/network"
)

type Conn struct {
	Connection *net.UDPConn
	Addr       net.UDPAddr

	closed chan bool

	packets chan []byte
}

func (connection *Conn) Close(listener *Listener) {
	delete(listener.connections, connection.Addr.String())
	delete(listener.actions, connection.Addr.String())
	connection.SendUnconnect()
	connection.closed <- true
}

func (connection *Conn) CloseFromClient(listener *Listener) {
	delete(listener.connections, connection.Addr.String())
	delete(listener.actions, connection.Addr.String())
	connection.closed <- true
}

func (connection *Conn) SendUnconnect() {
	b := bytes.NewBuffer(make([]byte, 0))
	(&network.Unconnected{}).Write(b)
	_, _ = connection.Connection.WriteToUDP(b.Bytes(), &connection.Addr)
}

func (connection *Conn) Write(bytes []byte) (n int, err error) {
	return connection.Connection.WriteToUDP(bytes, &connection.Addr)
}

func (connection *Conn) Read() (b []byte, err error) {
	select {
	case <-connection.closed:
		return nil, fmt.Errorf("connection closed %x", connection.Addr)
	case packet := <-connection.packets:
		return packet, nil
	}
}
