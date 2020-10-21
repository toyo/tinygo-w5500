package w5500

import (
	"errors"

	"github.com/toyo/tinygo-w5500/wiznet"
)

// SetDestIPPort set destination IP Port pair.
func (d *W5500) SetDestIPPort(socket uint8, ip wiznet.IP, port uint16) (err error) {
	if len(ip) != 4 {
		return errors.New(`Strange IP Address`)
	}
	err = d.WriteSocketRegisters(socket, 0x000c, ip) // Sn_DIPR (Socket n Destination IP Address Register)
	if err != nil {
		return
	}
	err = d.WriteSocketRegisters(socket, 0x0010, []byte{byte(port >> 8), byte(port & 0xff)}) // Sn_DPORT (Socket n Destination Port Register)
	return
}

// SocketOPEN sets socket OPEN
func (d *W5500) SocketOPEN(socket uint8) error {
	return d.WriteSocketRegisters(socket, 0x0001, []byte{0x01}) // OPEN
}

// SocketCLOSE sets socket CLOSE
func (d *W5500) SocketCLOSE(socket uint8) error {
	return d.WriteSocketRegisters(socket, 0x0001, []byte{0x10}) // CLOSE
}

// SocketLISTEN sets socket LISTEN
func (d *W5500) SocketLISTEN(socket uint8) error {
	return d.WriteSocketRegisters(socket, 0x0001, []byte{0x02}) // LISTEN
}

// SocketCONNECT sets socket CONNECT
func (d *W5500) SocketCONNECT(socket uint8) error {
	return d.WriteSocketRegisters(socket, 0x0001, []byte{0x04}) // CONNECT
}

// SocketDISCON sets socket DISCON
func (d *W5500) SocketDISCON(socket uint8) error {
	return d.WriteSocketRegisters(socket, 0x0001, []byte{0x08}) // DISCON
}
