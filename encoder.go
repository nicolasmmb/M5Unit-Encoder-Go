// Package encoder drives the M5Stack Unit Encoder (SKU U135), a rotary
// encoder with a push button and two SK6812 LEDs on I2C address 0x40.
package encoder

import (
	"fmt"
	"sync"
)

// DefaultAddress is the I2C address of the Unit Encoder out of the box.
const DefaultAddress = 0x40

// Registers exposed by the U135 firmware.
const (
	RegMode    = 0x00 // 1 byte: 0 = pulse, 1 = AB
	RegEncoder = 0x10 // int16 little-endian (read/write)
	RegButton  = 0x20 // 1 byte: 0 = pressed (active low)
	RegLED     = 0x30 // [index, R, G, B]
	RegReset   = 0x40 // write 1 to reset the counter
)

// Bus is the I2C transport. machine.I2C satisfies it through NewI2C.
type Bus interface {
	Tx(addr uint16, w, r []byte) error
}

// Device is a Unit Encoder on an I2C bus.
//
// A Device is safe for concurrent use.
type Device struct {
	mu   sync.Mutex
	bus  Bus
	addr uint16
}

// New returns a Device on bus. A zero addr uses DefaultAddress.
func New(bus Bus, addr uint16) *Device {
	if addr == 0 {
		addr = DefaultAddress
	}
	return &Device{bus: bus, addr: addr}
}

// Address returns the I2C address in use.
func (d *Device) Address() uint16 { return d.addr }

func (d *Device) read(reg byte, n int) ([]byte, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	buf := make([]byte, n)
	if err := d.bus.Tx(d.addr, []byte{reg}, buf); err != nil {
		return nil, fmt.Errorf("encoder: read 0x%02X: %w", reg, err)
	}
	return buf, nil
}

func (d *Device) write(reg byte, data ...byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	payload := make([]byte, 1+len(data))
	payload[0] = reg
	copy(payload[1:], data)
	if err := d.bus.Tx(d.addr, payload, nil); err != nil {
		return fmt.Errorf("encoder: write 0x%02X: %w", reg, err)
	}
	return nil
}

// Init checks that the module answers on the bus.
func (d *Device) Init() error {
	_, err := d.Encoder()
	return err
}

// Encoder returns the signed 16-bit encoder count (register 0x10).
func (d *Device) Encoder() (int16, error) {
	data, err := d.read(RegEncoder, 2)
	if err != nil {
		return 0, err
	}
	return int16(uint16(data[0]) | uint16(data[1])<<8), nil
}

// SetEncoder writes the encoder count (register 0x10).
func (d *Device) SetEncoder(value int16) error {
	u := uint16(value)
	return d.write(RegEncoder, byte(u), byte(u>>8))
}

// Button reports whether the encoder is pressed (register 0x20).
// The U135 firmware returns 0 while pressed and 1 while released.
func (d *Device) Button() (bool, error) {
	data, err := d.read(RegButton, 1)
	if err != nil {
		return false, err
	}
	return data[0] == 0, nil
}

// Reset zeroes the encoder counter (register 0x40, write 1).
func (d *Device) Reset() error {
	return d.write(RegReset, 1)
}

// Mode reads the work mode (register 0x00).
func (d *Device) Mode() (uint8, error) {
	data, err := d.read(RegMode, 1)
	if err != nil {
		return 0, err
	}
	return data[0], nil
}

// SetMode writes the work mode (register 0x00).
func (d *Device) SetMode(mode uint8) error {
	return d.write(RegMode, mode)
}

// SetLED writes the RGB LEDs. index 0 targets both LEDs; 1 and 2 target a
// single LED.
func (d *Device) SetLED(index, r, g, b uint8) error {
	if index > 2 {
		return fmt.Errorf("encoder: led index %d out of range (0..2)", index)
	}
	return d.write(RegLED, index, r, g, b)
}

// SetLEDColor writes one LED using a 0xRRGGBB color.
func (d *Device) SetLEDColor(index uint8, rgb uint32) error {
	return d.SetLED(index, byte(rgb>>16), byte(rgb>>8), byte(rgb))
}
