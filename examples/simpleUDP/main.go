package main

import (
	"fmt"
	"machine"
	"runtime"
	"time"

	"github.com/toyo/tinygo-w5500/w5500"
	"github.com/toyo/tinygo-w5500/wiznet"
)

func main() {

	machine.SPI0.Configure(machine.SPIConfig{})

	d, err := w5500.NewW5500(machine.SPI0, machine.D3,
		wiznet.IP{192, 168, 50, 55}, wiznet.IP{255, 255, 255, 0}, wiznet.IP{192, 168, 50, 1},
		[]byte{1, 2, 3, 4, 5, 6})

	time.Sleep(3 * time.Second)
	if err != nil {
		fmt.Println(err)
		time.Sleep(24 * 60 * 60 * time.Second)
		return
	}

	fmt.Println(d.ShowStatus()) // show like `FullDup 100Mbps LinkUp`

	udpTerminalClient(d, &wiznet.UDPAddr{IP: wiznet.IP{192, 168, 50, 3}, Port: 7777}) // 7 is echo

	time.Sleep(24 * 60 * 60 * time.Second)
}

func udpTerminalClient(d wiznet.Chip, raddr *wiznet.UDPAddr) {

	fmt.Println(`Client Dialing`)
	c, err := d.DialUDP(`udp`, nil, raddr)
	if err == nil {
		fmt.Println(`Client Connected`)

		for {
			b := make([]byte, 128)
			leng, peer, err := c.ReadFromUDP(b)
			if err != nil {
				fmt.Println(err)
				break
			}
			if leng != 0 {
				fmt.Println(`"`, string(b[:leng]), `" from`, peer)
			}

			char, ok := machine.UART0.Buffer.Get()
			if ok {
				n, err := c.WriteToUDP([]byte{'-', char, '-'}, raddr)
				if err != nil {
					fmt.Println(`WriteToUDP error`, err, n)
				}
			}

			if c.IsCLOSED() {
				fmt.Println(`Client Closed`)
				break
			}
			runtime.GC()
		}
	} else {
		fmt.Println(err)
	}
	time.Sleep(24 * 60 * 60 * time.Second)
}
