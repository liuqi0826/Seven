package network

import (
	"net"
	"time"
)

const (
	NETWORK_TYPE_LOCAL      = "local"
	NETWORK_TYPE_TCP        = "tcp"
	NETWORK_TYPE_UDP        = "udp"
	NETWORK_TYPE_WEB_SOCKET = "webSocket"
)

func TCPConnect(n *Network, addr string) {
	for {
		conn, err := net.Dial("tcp", addr)
		if err == nil {
			n.Create(conn)
			break
		} else {
			time.Sleep(time.Microsecond * 10)
		}
	}
}

func LocalConnect(n *Network, l *Listener) {
	c := new(Connection)
	c.Connection()
	n.Create(c)
	l.localChan <- c
}
