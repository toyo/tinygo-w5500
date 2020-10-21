package w5500

import "strconv"

func (d *W5500) getStatus(socket uint8) (uint8, error) {
	x, err := d.ReadSocketRegisters(socket, 0x0003, 1)
	return x[0], err
}

// GetSocketStatus returns socket status in string.
func (d *W5500) GetSocketStatus(socket uint8) (s string) {
	x, err := d.ReadSocketRegisters(socket, 0x0003, 1)
	if err != nil {
		return err.Error()
	}
	switch x[0] {
	case 0x00:
		return `SOCK_CLOSED`
	case 0x13:
		return `SOCK_INIT`
	case 0x14:
		return `SOCK_LISTEN`
	case 0x17:
		return `SOCK_ESTABLISHED`
	case 0x1c:
		return `SOCK_CLOSE_WAIT`
	case 0x22:
		return `SOCK_UDP`
	case 0x42:
		return `SOCK_MACRAW`
	case 0x15:
		return `SOCK_SYNSENT`
	case 0x16:
		return `SOCK_SYNRECV`
	case 0x18:
		return `SOCK_FIN_WAIT`
	case 0x1a:
		return `SOCK_CLOSING`
	case 0x1b:
		return `SOCK_TIME_WAIT`
	case 0x1d:
		return `SOCK_LAST_ACK`
	default:
		return `SOCK_???` + strconv.FormatUint(uint64(x[0]), 10)
	}
}

// IsSocketCLOSED returns the status is CLOSED or not.
func (d *W5500) IsSocketCLOSED(socket uint8) bool {
	reg, err := d.getStatus(socket)
	return reg == 0x00 && err == nil
}

// IsSocketESTABLISHED return true if tcp status is ESTABLISHED
func (d *W5500) IsSocketESTABLISHED(socket uint8) bool {
	status, err := d.getStatus(socket)
	return status == 0x17 && err == nil
}

// IsSocketCLOSEWAIT returns the status is CLOSE_WAIT or not.
func (d *W5500) IsSocketCLOSEWAIT(socket uint8) bool {
	status, err := d.getStatus(socket)
	return status == 0x1c && err == nil
}
