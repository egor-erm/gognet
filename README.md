# GOGNET

Gognet is a library for creating simple game servers based on the udp protocol. At its core, it represents a simplified version of the raknet gaming protocol, but without unnecessary functions. This allows you to create simpler clients in different programming languages to work with gognet. To begin with, the client establishes a connection with the server, then can safely exchange data with it.

# Getting Started
To create a server:
```go
package main

import (
	"net"

	"github.com/egor-erm/gognet"
)

func main() {
	listener, _ := gognet.Listen(net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 15000}) // ip and port of listener
		
	defer listener.Close()

	for {
		conn, _ := listener.Accept() //new connection
    
		go func() { //create gorutine for connection
      			defer conn.Close()
			for {
				b, _ := conn.Read() // read []byte packet

				conn.Write(b) // write []byte packet
			}
		}()
	}
}
```

To create a client:

```go
package main

import (
	"bytes"
	"fmt"
	"net"

	"github.com/egor-erm/gognet"
)

var remoteAdr = net.UDPAddr{IP: []byte{0, 0, 0, 0}, Port: 15000}

func main() {
	conn, _ := gognet.Dial(&remoteAdr) //crate client
	defer conn.Close()

	for {
		_, _ = conn.Write([]byte{1, 2, 3, 4, 5}) // write []byte packet

		b := make([]byte, 1024*1024*3)
		n, _ := conn.Read(b) // read []byte packet

		buf := bytes.NewBuffer(b[:n])

		fmt.Println(buf.Bytes()) // print packet
	}
}
```

## Telegram - https://t.me/ermolaymc
