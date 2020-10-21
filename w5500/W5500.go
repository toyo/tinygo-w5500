package w5500

import (
	"errors"
	"machine"
	"runtime"
	"sync"

	"github.com/toyo/tinygo-w5500/wiznet"
)

// W5500 is device class for W5500 Ethernet Chip.
type W5500 struct {
	spi              machine.SPI
	csPin            machine.Pin
	socketInUseMutex []sync.Mutex
}

// NewW5500 is constructor for W5500
func NewW5500(spi machine.SPI, csPin machine.Pin, myip, subnetmask, gatewayip wiznet.IP, macaddress []byte) (d *W5500, err error) {
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	csPin.High()
	d = &W5500{
		spi:              spi,
		csPin:            csPin,
		socketInUseMutex: make([]sync.Mutex, 8),
	}

	if err = d.writeCommonRegister(0x0000, 0x80); err != nil { // RST (Reset) is Register0 bit 7
		return nil, err
	}

	if len(macaddress) != 6 {
		return nil, errors.New(`Wrong MAC Address`)
	}

	var initialaddress []byte
	initialaddress = append(initialaddress, []byte(gatewayip.To4())...)
	initialaddress = append(initialaddress, []byte(subnetmask.To4())...)
	initialaddress = append(initialaddress, macaddress...)
	initialaddress = append(initialaddress, []byte(myip.To4())...)
	if err = d.writeCommonRegisters(0x0001, initialaddress); err != nil {
		return nil, err
	}

	return d, nil
}

// DialTCP is for TCP client connection.
func (d *W5500) DialTCP(network string, laddr, raddr *wiznet.TCPAddr) (*wiznet.TCPConn, error) {
	socket, err := d.TCPSocketOpen(0)
	if err != nil {
		return nil, err
	}
	d.SetDestIPPort(socket, raddr.IP.To4(), uint16(raddr.Port))
	d.SocketCONNECT(socket)
	for !d.IsSocketESTABLISHED(socket) { // wait SOCK_ESTABLISHED
		runtime.Gosched()
	}

	return &wiznet.TCPConn{Socket: wiznet.NewSocket(d, socket)}, nil
}

// ListenTCP returns TCPListener
func (d *W5500) ListenTCP(network string, laddr *wiznet.TCPAddr) (*wiznet.TCPListener, error) {
	return wiznet.NewTCPListener(d, uint16(laddr.Port))
}

// DialUDP is for UDP client connection.
func (d *W5500) DialUDP(network string, laddr, raddr *wiznet.UDPAddr) (*wiznet.UDPConn, error) {
	socket, err := d.udpSocketOpen(0)
	if err != nil {
		return nil, err
	}
	return &wiznet.UDPConn{Socket: wiznet.NewSocket(d, socket)},
		d.SetDestIPPort(socket, raddr.IP.To4(), uint16(raddr.Port))
}

// ListenUDP returns UDPListener
func (d *W5500) ListenUDP(network string, laddr *wiznet.UDPAddr) (*wiznet.UDPConn, error) {
	socket, err := d.udpSocketOpen(uint16(laddr.Port))
	return &wiznet.UDPConn{Socket: wiznet.NewSocket(d, socket)}, err
}
