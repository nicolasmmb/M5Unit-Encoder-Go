# M5Unit-Encoder-Go

[![test](https://github.com/nicolasmmb/M5Unit-Encoder-Go/actions/workflows/test.yml/badge.svg)](https://github.com/nicolasmmb/M5Unit-Encoder-Go/actions/workflows/test.yml)

Pure-Go driver for the **M5Stack Unit Encoder (SKU U135)**: a rotary encoder
with a push button and two SK6812 RGB LEDs on I2C.

## Features

- Encoder count (signed 16-bit), button state, work mode and both RGB LEDs.
- Small API with an injectable I2C `Bus`; the core has no hardware imports and
  is fully testable on the host.
- Safe for concurrent use.
- No C, C++ or CGo.

## Hardware

| Parameter | Value |
|---|---|
| MCU | STM32F030F4P6 |
| I2C address | `0x40` |
| Encoder | 30-position pulse knob with push button |
| LEDs | 2× SK6812 |
| Wiring (M5StickC) | SDA = G32, SCL = G33, Grove I2C1 |

The Grove port is shared with other I2C units such as ENV IV and the
Fingerprint2; only one can be connected at a time.

## Install

```bash
go get github.com/nicolasmmb/M5Unit-Encoder-Go@v0.1.0
```

## Quick start

```go
i2c := machine.I2C1
_ = i2c.Configure(machine.I2CConfig{
    SDA:       machine.GPIO32,
    SCL:       machine.GPIO33,
    Frequency: 100 * machine.KHz,
})

dev := encoder.NewI2C(i2c, encoder.DefaultAddress)
if err := dev.Init(); err != nil {
    panic(err)
}

count, _ := dev.Encoder()   // signed 16-bit
pressed, _ := dev.Button()
_ = dev.SetLEDColor(0, 0xFF0000)
_ = dev.SetEncoder(0)
```

## Registers

| Register | Name | Content |
|---|---|---|
| `0x00` | Mode (write-only) | 1 byte: 0 = pulse, 1 = AB |
| `0x10` | Encoder | `int16` little-endian |
| `0x20` | Button | 1 byte, 0 = pressed (active low) |
| `0x30` | RGB LED | `[index, R, G, B]`, index 0 = both, 1/2 = one LED |
| `0x40` | Reset | write 1 to zero the counter |

## API

| Function | Description |
|---|---|
| `New(bus, addr)` | Device on a custom `Bus` (0 uses `0x40`) |
| `NewI2C(i2c, addr)` | Device on a TinyGo `machine.I2C` (esp32) |
| `Init()` | Probe the module |
| `Encoder()` / `SetEncoder(v)` | Read/write the count |
| `Button()` | Button state |
| `SetMode(m)` | Work mode: `ModePulse` (0) or `ModeAB` (1); write-only |
| `SetLED(index, r, g, b)` | Set LEDs (0 = both, 1/2 = one) |
| `SetLEDColor(index, rgb)` | Set LEDs with `0xRRGGBB` |
| `Reset()` | Zero the encoder counter |

## References

- [M5Stack Unit Encoder (U135)](https://docs.m5stack.com/en/unit/encoder)
- [M5Stack Arduino library](https://github.com/m5stack/UNIT_ENCODER)
- [I2C protocol PDF](https://m5stack-doc.oss-cn-shenzhen.aliyuncs.com/801/Unit-Encoder_Protocol.pdf)

## Testing

```bash
make test     # go test ./...
make examples # build the esp32 samples
```

## License

MIT — see [LICENSE](LICENSE).
