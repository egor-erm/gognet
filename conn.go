package gognet

import (
	"net"
	"strconv"
	"strings"
)

type Conn struct {
	conn net.UDPConn
}

func Convert(addr string) *net.UDPAddr {
	ip := strings.Split(addr, ":")[0]
	port, _ := strconv.Atoi(strings.Split(addr, ":")[1])

	return &net.UDPAddr{IP: net.IP(ip), Port: port}
}
