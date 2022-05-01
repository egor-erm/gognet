package gognet

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"sync"

	"github.com/egor-erm/gognet/network"
)

type Listener struct {
	listener *net.UDPConn

	incoming chan *net.Conn
	closed   chan *net.Conn

	connections sync.Map

	log        *log.Logger
	listenerId int32
}

var listenerID = rand.Int31()

func Listen(address *net.UDPAddr) (*Listener, error) {
	list, err := net.ListenUDP("udp", address)
	if err != nil {
		return nil, &net.OpError{Op: "listen", Net: "gognet", Source: nil, Addr: nil, Err: err}
	}
	listener := &Listener{
		listener: list,

		incoming: make(chan *net.Conn),
		closed:   make(chan *net.Conn),

		log:        log.New(os.Stderr, "", log.LstdFlags),
		listenerId: listenerID,
	}

	go listener.listen()
	return listener, nil
}

func (listener *Listener) listen() {
	b := make([]byte, 1500)
	buf := bytes.NewBuffer(b[:0])
	for {
		n, addr, err := listener.listener.ReadFromUDP(b)
		if err != nil {
			panic(err)
			return
		}
		_, _ = buf.Write(b[:n])

		if err := listener.handle(buf, addr); err != nil {
			listener.log.Printf("listener: error handling packet (addr = %v): %v\n", addr, err)
		}
		buf.Reset()
	}
}

func (listener *Listener) handle(b *bytes.Buffer, addr net.Addr) error {
	_, found := listener.connections.Load(addr.String())
	if !found {
		// If there was no session yet, it means the packet is an offline message. It is not contained in a
		// datagram.
		packetID, err := b.ReadByte()
		if err != nil {
			return fmt.Errorf("error reading packet ID byte: %v", err)
		}
		switch packetID {
		case network.IDOpenConnectionRequest1:
			return listener.handleOpenConnectionRequest1(b, addr)
		default:
			return fmt.Errorf("unknown packet received (%x): %x", packetID, b.Bytes())
		}
		return nil
	}
	return nil
}

func (listener *Listener) handleOpenConnectionRequest1(b *bytes.Buffer, addr net.Addr) error {
	packet := &network.OpenConnectionRequest1{}
	if err := packet.Read(b); err != nil {
		return fmt.Errorf("error reading open connection request 1: %v", err)
	}
	b.Reset()

	if packet.Protocol != network.Protocol_Version {
		(&network.IncompatibleProtocolVersion{ServerGUID: listener.listenerId, ServerProtocol: network.Protocol_Version}).Write(b)
		_, _ = listener.listener.WriteTo(b.Bytes(), addr)
		return fmt.Errorf("error handling open connection request 1: incompatible protocol version %v (listener protocol = %v)", packet.Protocol, network.Protocol_Version)
	}

	(&network.OpenConnectionReply1{ServerGUID: listener.listenerId}).Write(b)
	_, err := listener.listener.WriteTo(b.Bytes(), addr)
	return err
}
