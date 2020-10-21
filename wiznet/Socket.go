package wiznet

import (
	"fmt"
)

// Socket is for each socket.
type Socket struct {
	d      Chip
	socket uint8
}

// NewSocket creates NewSocket class.
func NewSocket(d Chip, socket uint8) *Socket {
	return &Socket{d: d, socket: socket}
}

// GetSocketNumber returns socket number.
func (c Socket) GetSocketNumber() uint8 {
	return c.socket
}

// String shows the status of the socket.
func (c Socket) String() (s string) {
	return fmt.Sprintf("Socket: %d, Status: 0x%02X\n", c.socket, c.d.GetSocketStatus(c.socket))
}

// Close close the connection.
func (c Socket) Close() error {
	if !c.d.IsSocketCLOSED(c.socket) { // SOCK_CLOSED
		return c.d.SocketCLOSE(c.socket)
	}
	return nil
}

// IsCLOSED is for outside
func (c Socket) IsCLOSED() bool {
	return c.d.IsSocketCLOSED(c.socket)
}
