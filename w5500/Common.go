package w5500

import (
	"errors"
	"math"

	"github.com/toyo/tinygo-w5500/wiznet"
)

// String shows the status of each socket.
func (d *W5500) String() (s string) {
	for i := uint8(0); i < uint8(len(d.socketInUseMutex)); i++ {
		s += wiznet.NewSocket(d, i).String() + "\n"
	}
	return
}

// ShowStatus shows ths status of Hardware.
func (d W5500) ShowStatus() (s string) {
	PHYConfig, err := d.readCommonRegister(0x002e) // GetPHYCFGR get W5500 PHY Configuration Register
	if err == nil {
		if PHYConfig&0b00000100 != 0 {
			s += `FullDup `
		} else {
			s += `HalfDup `
		}
		if PHYConfig&0b00000010 != 0 {
			s += `100Mbps `
		} else {
			s += ` 10Mbps `
		}
		if PHYConfig&0b00000001 != 0 {
			s += `LinkUp `
		} else {
			s += `Link Down `
		}
	} else {
		s += err.Error()
	}
	return
}

// TCPSocketOpen is for TCP socket open.
func (d *W5500) TCPSocketOpen(port uint16) (uint8, error) {
	return d.socketOpen(1, port)
}

// udpSocketOpen is for UDP socket open.
func (d *W5500) udpSocketOpen(port uint16) (uint8, error) {
	return d.socketOpen(2, port)
}

// Please set port 0 for TCP client.
// protocol: tcp=1, udp=2, MACraw=4(not implement)
func (d *W5500) socketOpen(protocol uint8, port uint16) (socket uint8, err error) {

	for socket = 0; socket < uint8(len(d.socketInUseMutex)); socket++ {
		d.socketInUseMutex[socket].Lock()
		if d.IsSocketCLOSED(socket) { // SOCK_CLOSED
			break
		}
		d.socketInUseMutex[socket].Unlock()
	}
	if socket == uint8(len(d.socketInUseMutex)) {
		err = errors.New(`No available socket`)
	}

	defer d.socketInUseMutex[socket].Unlock()

	if err != nil {
		return math.MaxUint8, err
	}

	if port == 0 {
		port = 1024 + uint16(socket) // Set LocalPort automagically.
	}

	if err = d.WriteSocketRegisters(socket, 0x0000, []byte{protocol & 0x0f}); err != nil { // Set TCP/UDP
		return math.MaxUint8, err
	}
	if err = d.writeSocketRegister16(socket, 0x0004, port); err != nil { // Set LocalPort
		return math.MaxUint8, err
	}
	if err = d.SocketOPEN(socket); err != nil { // Set OPEN
		return math.MaxUint8, err
	}
	if d.IsSocketCLOSED(socket) {
		return math.MaxUint8, errors.New(`Cannot INIT Socket`)
	}

	return
}
