package wiznet

import (
	"io"
)

// TCPConn is for each socket.
type TCPConn struct {
	*Socket
}

// ReadFrom is implement for io.ReaderFrom
func (c TCPConn) ReadFrom(r io.Reader) (n int64, err error) {
	buf := make([]byte, 146)
	var nn int
	nn, err = r.Read(buf)
	if err == nil && nn != 0 {
		nn, err = c.Write(buf[:nn])
		n = int64(nn)
	}
	return
}

// WriteTo is implement for io.WriterTo
func (c TCPConn) WriteTo(w io.Writer) (n int64, err error) {
	buf := make([]byte, 146)
	var nn int
	nn, err = c.Read(buf)
	if err == nil && nn != 0 {
		nn, err = w.Write(buf[:nn])
		n = int64(nn)
	}
	return
}

func (c TCPConn) Read(p []byte) (int, error) {
	length := uint16(len(p))
	if length != 0 {
		data, err := c.d.SocketRead(c.Socket.socket, length)
		if err != nil {
			return 0, err
		}
		copy(p, data)
		if uint16(len(p)) > length {
			p[length] = 0x00 // terminate for string
		}

	}
	return int(length), nil
}

func (c TCPConn) Write(p []byte) (int, error) {
	if err := c.d.SocketWrite(c.Socket.socket, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

// Close closes the connection.
func (c TCPConn) Close() error {
	if !c.d.IsSocketCLOSED(c.socket) { // SOCK_CLOSED
		return c.d.SocketDISCON(c.socket)
	}

	return nil
}

// IsCLOSEWAIT returns the status is CLOSE_WAIT or not.
func (c TCPConn) IsCLOSEWAIT() bool {
	return c.d.IsSocketCLOSEWAIT(c.socket)
}
