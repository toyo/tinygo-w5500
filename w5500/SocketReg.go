package w5500

import "encoding/binary"

// ReadSocketRegisters reads socket register.
func (d *W5500) ReadSocketRegisters(socket uint8, address uint16, size uint16) (data []byte, err error) {
	return d.spiread(address, socket<<2|0b0001, size)
}

// WriteSocketRegisters writes socket register.
func (d *W5500) WriteSocketRegisters(socket uint8, address uint16, data []byte) (err error) {
	return d.spiwrite(address, socket<<2|0b00001, data)
}

func (d *W5500) writeSocketRegister16(socket uint8, address uint16, data uint16) (err error) {
	databyte := make([]byte, 2)
	binary.BigEndian.PutUint16(databyte, data)
	return d.WriteSocketRegisters(socket, address, databyte)
}
