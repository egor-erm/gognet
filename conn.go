package gognet

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

type Conn struct {
	conn net.UDPConn
	addr net.UDPAddr

	closed chan struct{}

	packets chan []byte
}

func Convert(addr string) *net.UDPAddr {
	ip := strings.Split(addr, ":")[0]
	port, _ := strconv.Atoi(strings.Split(addr, ":")[1])

	return &net.UDPAddr{IP: net.IP(ip), Port: port}
}

func (connection *Conn) Write(bytes *[]byte) (n int, err error) {
	return connection.conn.WriteToUDP(*bytes, &connection.addr)
}

func (connection *Conn) Read(bytes *[]byte) (b []byte, err error) {
	select {
	case packet := <-connection.packets:
		return packet, err
	case <-connection.closed:
		return nil, fmt.Errorf("connection closed %x", connection.addr)
	}
}
