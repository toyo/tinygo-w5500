package wiznet

import (
	"encoding/binary"
)

// UDPConn is for each socket.
type UDPConn struct {
	*Socket
	peerAddr   *UDPAddr
	remainSize uint16
}

// ReadFrom is same as net.UDPConn.ReadFrom()
func (c *UDPConn) ReadFrom(b []byte) (int, Addr, error) {
	return c.ReadFromUDP(b)
}

// ReadFromUDP is same as net.UDPConn.ReadFromUDP()
func (c *UDPConn) ReadFromUDP(b []byte) (leng int, peer *UDPAddr, err error) {
	length := uint16(len(b))
	if length != 0 {

		if c.remainSize == 0 { // read header
			var data []byte
			data, err = c.d.SocketRead(c.Socket.socket, 8)
			if err != nil {
				return
			}
			if len(data) == 0 {
				return 0, nil, nil
			}
			c.peerAddr = &UDPAddr{IP: data[:4], Port: int(binary.BigEndian.Uint16(data[4:6]))}
			c.remainSize = binary.BigEndian.Uint16(data[6:8])
			if c.remainSize < length {
				length = c.remainSize
			}
		}
		data, err := c.d.SocketRead(c.Socket.socket, length)
		if err != nil {
			return 0, nil, err
		}
		if len(data) == 0 {
			return 0, nil, nil
		}
		copy(b, data)
		if len(b) > int(length) {
			b[length] = 0x00 // terminate for string
		}
		c.remainSize -= length

	}
	leng = int(length)
	peer = c.peerAddr
	return
}

// WriteTo is same as net.UDPConn.WriteTo()
func (c *UDPConn) WriteTo(b []byte, addr Addr) (int, error) {
	return c.WriteToUDP(b, addr.(*UDPAddr))
}

// WriteToUDP is same as net.UDPConn.WriteToUDP()
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (n int, err error) {
	c.d.SetDestIPPort(c.socket, addr.IP.To4(), uint16(addr.Port))
	if err = c.d.SocketWrite(c.Socket.socket, b); err != nil {
		return
	}
	n = len(b)
	return
}
