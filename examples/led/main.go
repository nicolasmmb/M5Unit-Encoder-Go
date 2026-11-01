//go:build esp32

// Example: cycle the two LEDs through colors.
package main

import (
	"time"

	"github.com/nicolasmmb/M5Unit-Encoder-Go/examples/common"
)

func main() {
	time.Sleep(2 * time.Second)
	dev := common.Open()
	colors := []uint32{0x000000, 0xFF0000, 0x00FF00, 0x0000FF, 0xFFFFFF}
	for i := 0; ; i++ {
		c := colors[i%len(colors)]
		_ = dev.SetLEDColor(0, c)
		_ = dev.SetLEDColor(1, c)
		println("led color", c)
		time.Sleep(700 * time.Millisecond)
	}
}
