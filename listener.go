package gognet

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/egor-erm/gognet/network"
)

type Listener struct {
	listener *net.UDPConn

	incoming chan *net.UDPConn
	closed   chan *net.UDPConn

	connections sync.Map

	log        *log.Logger
	listenerId int32
}

var listenerID = rand.Int31()

func Listen(address string) (*Listener, error) {
	ip := strings.Split(address, ":")[0]
	port := strings.Split(address, ":")[1]
	p, _ := strconv.Atoi(port)

	la := &net.UDPAddr{IP: net.IP(ip), Port: p}

	list, err := net.ListenUDP("udp", la)
	if err != nil {
		return nil, &net.OpError{Op: "listen", Net: "gognet", Source: nil, Addr: nil, Err: err}
	}
	listener := &Listener{
		listener: list,

		incoming: make(chan *net.UDPConn),
		closed:   make(chan *net.UDPConn),

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
		}
		_, _ = buf.Write(b[:n])

		if err := listener.handle(buf, addr); err != nil {
			listener.log.Printf("listener: error handling packet (addr = %v): %v\n", addr, err)
		}
		buf.Reset()
	}
}

func (listener *Listener) Accept() (*net.UDPConn, error) {
	conn, ok := <-listener.incoming
	if !ok {
		return nil, &net.OpError{Op: "accept", Net: "raknet", Source: nil, Addr: nil, Err: fmt.Errorf("Conn closed")}
	}
	return conn, nil
}

func (listener *Listener) handle(b *bytes.Buffer, addr net.Addr) error {
	_, found := listener.connections.Load(addr.String())
	if !found {

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

	listener.connections.Store(addr, listener.listener)

	listener.incoming <- listener.listener
	return err
}
