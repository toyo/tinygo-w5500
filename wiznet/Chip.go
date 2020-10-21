package wiznet

// Chip is for WIZnet chip
type Chip interface {
	DialTCP(network string, laddr, raddr *TCPAddr) (*TCPConn, error)
	ListenTCP(network string, laddr *TCPAddr) (*TCPListener, error)
	DialUDP(network string, laddr, raddr *UDPAddr) (*UDPConn, error)
	ListenUDP(network string, laddr *UDPAddr) (*UDPConn, error)
	TCPSocketOpen(port uint16) (socket uint8, err error)

	SetDestIPPort(socket uint8, ip IP, port uint16) error
	SocketWrite(socket uint8, p []byte) error
	SocketRead(socket uint8, length uint16) (data []byte, err error)

	SocketLISTEN(socket uint8) error
	SocketDISCON(socket uint8) error
	SocketCLOSE(socket uint8) error

	IsSocketESTABLISHED(socket uint8) bool
	IsSocketCLOSEWAIT(socket uint8) bool
	IsSocketCLOSED(socket uint8) bool

	GetSocketStatus(socket uint8) (s string)
}
