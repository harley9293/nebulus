package tcpd

import "net"

type conn struct {
	net.Conn
}
