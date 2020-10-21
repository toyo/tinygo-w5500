package w5500

func (d W5500) readCommonRegisters(address uint16, size uint16) (data []byte, err error) {
	return d.spiread(address, 0, size)
}

func (d W5500) readCommonRegister(address uint16) (data byte, err error) {
	datas, err := d.readCommonRegisters(address, 1)
	return datas[0], err
}

func (d W5500) writeCommonRegisters(address uint16, data []byte) (err error) {
	return d.spiwrite(address, 0, data)
}

func (d W5500) writeCommonRegister(address uint16, data byte) (err error) {
	return d.writeCommonRegisters(address, []byte{data})
}
