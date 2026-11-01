package encoder

import (
	"errors"
	"sync"
	"testing"
)

type fakeBus struct {
	mu     sync.Mutex
	addr   uint16
	regs   map[byte][]byte
	writes [][]byte
	err    error
}

func newFakeBus() *fakeBus {
	return &fakeBus{addr: DefaultAddress, regs: map[byte][]byte{
		RegEncoder: {0x00, 0x00},
		RegButton:  {0x00},
		RegMode:    {0x00},
	}}
}

func (f *fakeBus) Tx(addr uint16, w, r []byte) error {
	f.mu.Lock()
	defer f.mu.Unlock()
	if f.err != nil {
		return f.err
	}
	if addr != f.addr {
		return errors.New("unexpected address")
	}
	if len(w) == 0 {
		return errors.New("missing register")
	}
	reg := w[0]
	if len(w) > 1 {
		f.regs[reg] = append([]byte(nil), w[1:]...)
		f.writes = append(f.writes, append([]byte{reg}, w[1:]...))
		return nil
	}
	data, ok := f.regs[reg]
	if !ok {
		return errors.New("unknown register")
	}
	copy(r, data)
	return nil
}

func TestEncoderSignedValue(t *testing.T) {
	bus := newFakeBus()
	bus.regs[RegEncoder] = []byte{0xFE, 0xFF} // -2
	dev := New(bus, 0)
	got, err := dev.Encoder()
	if err != nil {
		t.Fatal(err)
	}
	if got != -2 {
		t.Fatalf("Encoder = %d, want -2", got)
	}
}

func TestSetEncoder(t *testing.T) {
	bus := newFakeBus()
	dev := New(bus, 0)
	if err := dev.SetEncoder(0x1234); err != nil {
		t.Fatal(err)
	}
	want := []byte{RegEncoder, 0x34, 0x12}
	if len(bus.writes) != 1 || string(bus.writes[0]) != string(want) {
		t.Fatalf("writes = % X, want % X", bus.writes, want)
	}
}

func TestButton(t *testing.T) {
	bus := newFakeBus()
	dev := New(bus, 0)
	bus.regs[RegButton] = []byte{1}
	if ok, err := dev.Button(); err != nil || !ok {
		t.Fatalf("Button = %v, %v; want true, nil", ok, err)
	}
	bus.regs[RegButton] = []byte{0}
	if ok, err := dev.Button(); err != nil || ok {
		t.Fatalf("Button = %v, %v; want false, nil", ok, err)
	}
}

func TestMode(t *testing.T) {
	bus := newFakeBus()
	dev := New(bus, 0)
	bus.regs[RegMode] = []byte{2}
	if m, err := dev.Mode(); err != nil || m != 2 {
		t.Fatalf("Mode = %d, %v; want 2", m, err)
	}
	if err := dev.SetMode(1); err != nil {
		t.Fatal(err)
	}
	if got := bus.regs[RegMode][0]; got != 1 {
		t.Fatalf("mode register = %d, want 1", got)
	}
}

func TestLED(t *testing.T) {
	bus := newFakeBus()
	dev := New(bus, 0)
	if err := dev.SetLED(1, 0x12, 0x34, 0x56); err != nil {
		t.Fatal(err)
	}
	want := []byte{1, 0x12, 0x34, 0x56}
	if got := bus.regs[RegLED]; string(got) != string(want) {
		t.Fatalf("LED = % X, want % X", got, want)
	}
	if err := dev.SetLEDColor(0, 0xAABBCC); err != nil {
		t.Fatal(err)
	}
	want = []byte{0, 0xAA, 0xBB, 0xCC}
	if got := bus.regs[RegLED]; string(got) != string(want) {
		t.Fatalf("LED color = % X, want % X", got, want)
	}
	if err := dev.SetLED(2, 0, 0, 0); err == nil {
		t.Fatal("SetLED with index 2 returned nil")
	}
}

func TestBusError(t *testing.T) {
	bus := newFakeBus()
	bus.err = errors.New("nack")
	dev := New(bus, 0)
	if err := dev.Init(); err == nil {
		t.Fatal("Init returned nil for a failing bus")
	}
}

func TestConcurrentAccess(t *testing.T) {
	bus := newFakeBus()
	dev := New(bus, 0)
	done := make(chan error, 8)
	for i := 0; i < 8; i++ {
		go func() {
			for j := 0; j < 50; j++ {
				if _, err := dev.Encoder(); err != nil {
					done <- err
					return
				}
			}
			done <- nil
		}()
	}
	for i := 0; i < 8; i++ {
		if err := <-done; err != nil {
			t.Fatal(err)
		}
	}
}
