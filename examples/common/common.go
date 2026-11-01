//go:build esp32

// Package common contains setup helpers shared by the encoder examples.
package common

import (
	"machine"

	encoder "github.com/nicolasmmb/M5Unit-Encoder-Go"
)

// Open configures the Grove I2C bus (SDA=G32, SCL=G33) and waits for the
// encoder on its default address.
func Open() *encoder.Device {
	i2c := machine.I2C1
	_ = i2c.Configure(machine.I2CConfig{
		SDA:       machine.GPIO32,
		SCL:       machine.GPIO33,
		Frequency: 100 * machine.KHz,
	})
	dev := encoder.NewI2C(i2c, encoder.DefaultAddress)
	for dev.Init() != nil {
		println("encoder: waiting for module on 0x40")
	}
	return dev
}
