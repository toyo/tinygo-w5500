package w5500

import "encoding/binary"

// SocketWrite writes socket TX buffer
func (d *W5500) SocketWrite(socket uint8, p []byte) (err error) {
	if len(p) == 0 {
		return nil
	}
	var reg []byte
	if reg, err = d.ReadSocketRegisters(socket, 0x0020, 6); err == nil {
		sendMaxSize := binary.BigEndian.Uint16(reg[0:2])
		sendPtr := binary.BigEndian.Uint16(reg[4:6])
		if sendMaxSize < uint16(len(p)) {
			p = p[:sendMaxSize]
		}

		if err = d.spiwrite(sendPtr, socket<<2|0b10, p); err == nil {
			if err = d.writeSocketRegister16(socket, 0x0024, sendPtr+uint16(len(p))); err == nil { // Update Sn_TX_WR (Socket n TX Write Pointer Register)
				return d.WriteSocketRegisters(socket, 0x0001, []byte{0x20}) // Set SEND
			}
		}
	}
	return
}

// SocketRead reads socket RX buffer
func (d *W5500) SocketRead(socket uint8, length uint16) (data []byte, err error) {
	if length == 0 {
		return []byte{}, nil
	}
	var reg []byte
	reg, err = d.ReadSocketRegisters(socket, 0x0026, 4)
	if err == nil {
		recvdSize := binary.BigEndian.Uint16(reg[0:2])
		recvPtr := binary.BigEndian.Uint16(reg[2:4])
		if length > recvdSize {
			length = recvdSize
		}
		if data, err = d.spiread(recvPtr, socket<<2|0b00011, length); err == nil {
			if err = d.writeSocketRegister16(socket, 0x0028, recvPtr+length); err == nil {
				err = d.WriteSocketRegisters(socket, 0x0001, []byte{0x40}) // Set RECV
			}
		}
	}
	return
}
