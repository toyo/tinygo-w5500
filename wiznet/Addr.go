package wiznet

import "fmt"

// IP is IP Address. Same as net.IP . This slice should be 4byte.
type IP []byte

// Addr is same as net.Addr .
type Addr interface {
	Network() string // name of the network (for example, "tcp", "udp")
	String() string  // string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
}

// TCPAddr is same as net.TCPAddr
type TCPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone; added in Go 1.1
}

// UDPAddr is same as net.UDPAddr
type UDPAddr struct {
	IP   IP
	Port int
	Zone string // IPv6 scoped addressing zone; added in Go 1.1
}

// To4 is check for 4byte slice or not.
func (ip IP) To4() IP {
	if len(ip) != 4 {
		panic(`Unknown IP`)
	} else {
		return ip
	}
}

// Network is same as net.Addr.Network()
func (a *TCPAddr) Network() string {
	return `tcp`
}

// Network is same as net.Addr.Network()
func (a *UDPAddr) Network() string {
	return `udp`
}

// String is same as net.Addr.String()
func (a *TCPAddr) String() string {
	return fmt.Sprintf(`%d.%d.%d.%d:%d`, a.IP[0], a.IP[1], a.IP[2], a.IP[3], a.Port)
}

// String is same as net.Addr.String()
func (a *UDPAddr) String() string {
	return fmt.Sprintf(`%d.%d.%d.%d:%d`, a.IP[0], a.IP[1], a.IP[2], a.IP[3], a.Port)
}
