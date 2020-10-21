package w5500

func (d W5500) spiread(address uint16, bsb uint8, size uint16) (data []byte, err error) {

	size += 3 // for header
	wr := make([]byte, size)
	wr[0] = byte(address >> 8)
	wr[1] = byte(address & 0xff)
	wr[2] = bsb<<3 | 0b000 // read VDM
	re := make([]byte, size)

	d.csPin.Low()
	err = d.spi.Tx(wr, re)
	d.csPin.High()

	return re[3:], err
}

func (d W5500) spiwrite(address uint16, bsb uint8, data []byte) (err error) {
	wr := make([]byte, 3)
	wr[0] = byte(address >> 8)
	wr[1] = byte(address & 0xff)
	wr[2] = bsb<<3 | 0b100 // write VDM
	wr = append(wr, data...)

	d.csPin.Low()
	err = d.spi.Tx(wr, nil)
	d.csPin.High()

	return
}
