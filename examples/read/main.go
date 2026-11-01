//go:build esp32

// Example: print the encoder count and button state.
package main

import (
	"time"

	"github.com/nicolasmmb/M5Unit-Encoder-Go/examples/common"
)

func main() {
	time.Sleep(2 * time.Second)
	dev := common.Open()
	for {
		v, err := dev.Encoder()
		if err != nil {
			println("encoder:", err.Error())
			time.Sleep(time.Second)
			continue
		}
		btn, _ := dev.Button()
		mode, _ := dev.Mode()
		println("count", v, "button", btn, "mode", mode)
		time.Sleep(100 * time.Millisecond)
	}
}
