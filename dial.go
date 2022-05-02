package gognet

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/egor-erm/gognet/network"
)

func randPort() int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(2000) + 30000
}

func Dial(address *net.UDPAddr) (*net.UDPConn, error) {
	la := &net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: randPort()}
	ra := address

	udpConn, err := net.DialUDP("udp", la, ra)

	if err != nil {
		panic(err)
	}

	buf := bytes.NewBuffer(make([]byte, 0))

	packet := &network.OpenConnectionRequest1{Protocol: network.Protocol_Version}
	packet.Write(buf)

	_, err1 := udpConn.Write(buf.Bytes())

	if err1 != nil {
		panic(err1)
	}
	buf.Reset()

	b := make([]byte, 1500)
	n, err1 := udpConn.Read(b)

	if err1 != nil {
		panic(err1)
	}

	buf.Write(b[:n])

	byt, _ := buf.ReadByte()

	switch byt {
	case network.IDOpenConnectionReply1:
		return udpConn, nil
	case network.IDIncompatibleProtocolVersion:
		return nil, fmt.Errorf("protocol not supported")
	}

	return nil, fmt.Errorf("error open connection")
}
