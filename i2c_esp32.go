//go:build esp32

package encoder

import "machine"

// NewI2C returns a Device on a TinyGo I2C bus.
func NewI2C(i2c *machine.I2C, addr uint16) *Device {
	return New(i2c, addr)
}
