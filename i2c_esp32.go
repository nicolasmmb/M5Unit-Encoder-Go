//go:build esp32 || esp32c3 || esp32s3 || rp2040

package encoder

import "machine"

// NewI2C returns a Device on a TinyGo I2C bus.
//
// The helper lives behind a chip tag because it takes the concrete
// machine.I2C type. The tag covers the chips whose I2C API matches: ESP32,
// ESP32-C3, ESP32-S3 and RP2040. On any other chip use New(bus, addr).
func NewI2C(i2c *machine.I2C, addr uint16) *Device {
	return New(i2c, addr)
}
