package gognet

import (
	"bytes"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/egor-erm/gognet/network"
)

func Dial(address string) (*net.UDPConn, error) {
	ip := strings.Split(address, ":")[0]
	port := strings.Split(address, ":")[1]
	p, _ := strconv.Atoi(port)

	la := &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 54169}
	ra := &net.UDPAddr{IP: net.IP(ip), Port: p}

	udpConn, err := net.DialUDP("udp", la, ra)

	if err != nil {
		panic(err)
	}

	packet := &network.OpenConnectionRequest1{Protocol: network.Protocol_Version}
	buf := bytes.NewBuffer(make([]byte, 0))
	packet.Write(buf)
	udpConn.Write(buf.Bytes())

	b := make([]byte, 1500)
	n, _, err1 := udpConn.ReadFromUDP(b)

	if err1 != nil {
		panic(err1)
	}

	buf.Write(b[:n])

	byt, _ := buf.ReadByte()

	switch byt {
	case network.IDOpenConnectionReply1:
		fmt.Println("Успех")
		return udpConn, nil
	case network.IDIncompatibleProtocolVersion:
		fmt.Println("Протокол не тянет")
	}

	return nil, nil
}
