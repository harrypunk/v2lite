package proxy

import (
	"io"
	"net"
)

type Inbound interface {
	Name() string
	Addr() string
	Handshake(underlay net.Conn) (io.ReadWriter, error)
}
