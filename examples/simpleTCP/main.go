package main

import (
	"fmt"
	"io"
	"machine"
	"runtime"
	"time"

	"github.com/toyo/tinygo-w5500/w5500"
	"github.com/toyo/tinygo-w5500/wiznet"
)

func main() {

	machine.SPI0.Configure(machine.SPIConfig{})

	time.Sleep(3 * time.Second)

	d, err := w5500.NewW5500(machine.SPI0, machine.D3,
		wiznet.IP{192, 168, 50, 55}, wiznet.IP{255, 255, 255, 0}, wiznet.IP{192, 168, 50, 1},
		[]byte{1, 2, 3, 4, 5, 6})

	if err != nil {
		fmt.Println(err)
		time.Sleep(24 * 60 * 60 * time.Second)
		return
	}

	fmt.Println(d.ShowStatus()) // show like `FullDup 100Mbps LinkUp`

	go tcpEchoServer(d, 5555) // Run echo server on port 5555

	tcpTerminalClient(d, &wiznet.TCPAddr{IP: wiznet.IP{192, 168, 50, 3}, Port: 7}) // 7 is echo

	time.Sleep(24 * 60 * 60 * time.Second)
}

func tcpTerminalClient(d wiznet.Chip, raddr *wiznet.TCPAddr) {

	fmt.Println(`Client Dialing`)
	c, err := d.DialTCP(`tcp`, nil, raddr)
	if err == nil {
		fmt.Println(`Client Connected`)

		for {
			if _, err = io.Copy(machine.UART0, c); err != nil {
				fmt.Println(err)
				break
			}
			if _, err = io.Copy(c, machine.UART0); err != nil {
				fmt.Println(err)
				break
			}
			if c.IsCLOSED() {
				fmt.Println(`Client Closed`)
				break
			}
			runtime.Gosched()
		}
	} else {
		fmt.Println(err)
	}
}

func tcpEchoServer(d wiznet.Chip, port uint16) {
	l, err := d.ListenTCP(`tcp`, &wiznet.TCPAddr{Port: int(port)})
	if err == nil {
		for {
			c, err := l.AcceptTCP()
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println(`Server Connected`)

			go func(c *wiznet.TCPConn) {
				for {
					if _, err := io.Copy(c, c); err != nil { // echo server
						fmt.Println(err)
						c.Close()
					}
					if c.IsCLOSEWAIT() { // disconnected by peer.
						c.Close()
					}
					if c.IsCLOSED() {
						fmt.Println(`Server Closed`)
						break
					}
					runtime.Gosched()
				}
			}(c)
		}
	} else {
		fmt.Println(err)
	}
}
