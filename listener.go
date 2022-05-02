package gognet

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"sync"

	"github.com/egor-erm/gognet/network"
)

type Listener struct {
	listener *net.UDPConn

	incoming chan *Conn
	closed   chan net.UDPConn

	connections sync.Map

	listenerId int32
}

var listenerID = rand.Int31()

func Listen(address net.UDPAddr) (*Listener, error) {
	list, err := net.ListenUDP("udp", &address)
	if err != nil {
		return nil, err
	}
	listener := &Listener{
		listener: list,

		incoming: make(chan *Conn),
		closed:   make(chan net.UDPConn),

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

		if err := listener.handle(buf, *addr); err != nil {
			fmt.Printf("listener: error handling packet (addr = %v): %v\n", addr, err)
		}
		buf.Reset()
	}
}

func (listener *Listener) Accept() (*Conn, error) {
	conn, ok := <-listener.incoming
	if !ok {
		return nil, &net.OpError{Op: "accept", Net: "gognet", Source: nil, Addr: nil, Err: fmt.Errorf("conn closed")}
	}
	return conn, nil
}

func (listener *Listener) Close() error {
	err := listener.listener.Close()
	return err
}

func (listener *Listener) handle(b *bytes.Buffer, addr net.UDPAddr) error {
	value, found := listener.connections.Load(addr.String())
	if !found {
		packetID, err := b.ReadByte()
		if err != nil {
			return fmt.Errorf("error reading packet ID byte: %v", err)
		}
		switch packetID {
		case network.IDOpenConnectionRequest1:
			fmt.Println("res open packet")
			return listener.handleOpenConnectionRequest1(b, addr)
		default:
			return fmt.Errorf("unknown packet received (%x): %x", packetID, b.Bytes())
		}
	}

	conn := value.(*Conn)

	select {
	case <-conn.closed:
		listener.connections.Delete(addr.String())
	default:
		conn.packets <- b.Bytes()
	}

	return nil
}

func (listener *Listener) handleOpenConnectionRequest1(b *bytes.Buffer, addr net.UDPAddr) error {
	packet := &network.OpenConnectionRequest1{}
	if err := packet.Read(b); err != nil {
		return fmt.Errorf("error reading open connection request 1: %v", err)
	}

	b.Reset()

	if packet.Protocol != network.Protocol_Version {
		(&network.IncompatibleProtocolVersion{ServerGUID: listener.listenerId, ServerProtocol: network.Protocol_Version}).Write(b)
		_, _ = listener.listener.Write(b.Bytes())
		return fmt.Errorf("error handling open connection request 1: incompatible protocol version %v (listener protocol = %v)", packet.Protocol, network.Protocol_Version)
	}

	(&network.OpenConnectionReply1{ServerGUID: listener.listenerId}).Write(b)
	_, err := listener.listener.WriteToUDP(b.Bytes(), &addr)

	conn := &Conn{conn: listener.listener, addr: addr, packets: make(chan []byte)}

	listener.connections.Store(addr.String(), conn)

	listener.incoming <- conn
	return err
}
