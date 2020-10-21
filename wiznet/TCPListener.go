package wiznet

import (
	"math"
	"runtime"
)

// TCPListener listen TCP socket.
type TCPListener struct {
	chip   Chip
	ipport uint16
	socket [2]uint8
}

// NewTCPListener is constructor for TCPListener
func NewTCPListener(d Chip, port uint16) (*TCPListener, error) {
	l := &TCPListener{chip: d, ipport: port}
	var err error
	for i := 0; i < len(l.socket); i++ {
		var socket uint8
		if socket, err = d.TCPSocketOpen(port); err != nil {
			return l, err
		}
		if err = d.SocketLISTEN(socket); err != nil { // Set LISTEN
			return l, err
		}
		l.socket[i] = socket
	}
	return l, err
}

// AcceptTCP accepts TCP connection.
func (l *TCPListener) AcceptTCP() (c *TCPConn, err error) {
	for {
		for i := 0; i < len(l.socket); i++ {
			if l.socket[i] != math.MaxUint8 && l.chip.IsSocketESTABLISHED(l.socket[i]) { // wait SOCK_ESTABLISHED
				c = &TCPConn{Socket: NewSocket(l.chip, l.socket[i])}
				err = nil

				var socket uint8
				if socket, err = l.chip.TCPSocketOpen(l.ipport); err != nil {
					return
				}
				if err = l.chip.SocketLISTEN(socket); err != nil { // Set LISTEN
					return
				}
				l.socket[i] = socket

				return
			}
			runtime.Gosched()
		}
	}
}

// Close close the connection.
func (l TCPListener) Close() (err error) {
	for i := 0; i < len(l.socket); i++ {
		err = Socket{d: l.chip, socket: l.socket[i]}.Close()
		if err != nil {
			return
		}
	}
	return nil
}
