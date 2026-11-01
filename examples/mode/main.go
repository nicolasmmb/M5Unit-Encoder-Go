//go:build esp32

// Example: read and toggle the work mode (register 0x00).
package main

import (
	"time"

	"github.com/nicolasmmb/M5Unit-Encoder-Go/examples/common"
)

func main() {
	time.Sleep(2 * time.Second)
	dev := common.Open()
	m, err := dev.Mode()
	if err != nil {
		println("mode:", err.Error())
		return
	}
	println("current mode", m)
	next := uint8(0)
	if m == 0 {
		next = 1
	}
	if err := dev.SetMode(next); err != nil {
		println("set mode:", err.Error())
		return
	}
	println("mode set to", next)
	for {
		time.Sleep(time.Second)
	}
}
