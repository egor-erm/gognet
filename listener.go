package gognet

import (
	"bytes"
	"fmt"
	"math/rand"
	"net"
	"time"

	"github.com/egor-erm/gognet/network"
)

type Listener struct {
	listener *net.UDPConn

	incoming chan *Conn

	connections map[string]*Conn

	actions    map[string]int64
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

		connections: make(map[string]*Conn),
		actions:     make(map[string]int64),

		listenerId: listenerID,
	}

	go listener.listen()
	go listener.tick()
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
	conn, found := listener.connections[addr.String()]

	packetID, err := b.ReadByte()
	if !found {
		if err != nil {
			return fmt.Errorf("error reading packet ID byte: %v", err)
		}
		switch packetID {
		case network.IDOpenConnectionRequest1:
			return listener.handleOpenConnectionRequest1(b, addr)
		default:
			(&network.Unconnected{}).Write(b)
			_, _ = listener.listener.WriteToUDP(b.Bytes(), &addr)
			return fmt.Errorf("unknown unconnected packet received (%x): %x", packetID, b.Bytes())
		}
	}

	switch packetID {
	case network.IDUnconnected:
		listener.connections[addr.String()].Close(listener)
		return fmt.Errorf("unconnected packet received (%x): %x", packetID, b.Bytes())
	}

	listener.actions[addr.String()] = time.Now().Unix()
	conn.packets <- b.Bytes()

	return nil
}

func (listener *Listener) handleOpenConnectionRequest1(b *bytes.Buffer, addr net.UDPAddr) error {
	packet := &network.OpenConnectionRequest1{}
	if err := packet.Read(b); err != nil {
		return fmt.Errorf("error reading open connection request 1: %v", err)
	}

	b.Reset()

	if packet.Protocol != Protocol_Version {
		(&network.IncompatibleProtocolVersion{ServerGUID: listener.listenerId, ServerProtocol: Protocol_Version}).Write(b)
		_, _ = listener.listener.WriteToUDP(b.Bytes(), &addr)
		return fmt.Errorf("error handling open connection request 1: incompatible protocol version %v (listener protocol = %v)", packet.Protocol, Protocol_Version)
	}

	pack := &network.OpenConnectionReply1{ServerGUID: listener.listenerId}
	pack.Write(b)

	_, err := listener.listener.WriteToUDP(b.Bytes(), &addr)

	if err == nil {
		conn := &Conn{Connection: listener.listener, Addr: addr, packets: make(chan []byte), closed: make(chan bool)}

		listener.connections[addr.String()] = conn

		listener.actions[addr.String()] = time.Now().Unix()

		listener.incoming <- conn
	}

	return err
}

func (listener *Listener) tick() {
	for {
		time.Sleep(time.Second * 5)
		for key, value := range listener.actions {
			if time.Now().Unix() > value+timeout_time {
				listener.connections[key].Close(listener)
			}
		}
	}
}

func (listener *Listener) SetDeadLineTimeOut(sec int64) {
	timeout_time = sec
}
